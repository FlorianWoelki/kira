package pkg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/florianwoelki/kira/internal"
	"github.com/florianwoelki/kira/internal/cache"
	"github.com/florianwoelki/kira/internal/pool"
	"github.com/google/uuid"
)

const (
	// amountOfUsers is the amount of system users that are created and are available for
	// code execution.
	amountOfUsers = 50
	// maxOutputBufferCapacity is the maximum capacity of the output buffer.
	maxOutputBufferCapacity = "65332"
)

// RceEngine is the main struct which contains the worker pool, the system users
// and the cache.
type RceEngine struct {
	systemUsers *pool.SystemUsers
	pool        *pool.WorkerPool
	cache       *cache.Cache[pool.CodeOutput]
}

// NewRceEngine creates a new RceEngine instance that can be used to execute code.
func NewRceEngine() *RceEngine {
	return &RceEngine{
		systemUsers: pool.NewSystemUser(amountOfUsers),
		pool:        pool.NewWorkerPool(amountOfUsers),
		cache:       cache.NewCache[pool.CodeOutput](),
	}
}

func (rce *RceEngine) action(data pool.WorkData, ch chan<- pool.CodeOutput) {
	// Get the language by the name.
	language, err := GetLanguageByName(data.Lang)
	if err != nil {
		ch <- pool.CodeOutput{}
		return
	}

	// Check if the code is already in the cache, if so, return and send the cached output.
	var cacheOutput pool.CodeOutput
	if !data.BypassCache {
		cacheOutput, err = rce.cache.Get(language.Name, data.Code)

		if err == nil {
			ch <- cacheOutput
			return
		}
	}

	// Acquire a system user to execute the code.
	user, err := rce.systemUsers.Acquire()
	if err != nil {
		rce.systemUsers.Release(user.Uid)
		ch <- pool.CodeOutput{}
		return
	}

	tempDirName := uuid.New().String()

	// Create a temporary directory that is used to store the user's files in it.
	err = internal.CreateTempDir(user.Username, tempDirName)
	if err != nil {
		rce.systemUsers.Release(user.Uid)
		ch <- pool.CodeOutput{}
		return
	}

	// Create a temporary file that is used to store the user's code in it.
	filename, err := internal.CreateTempFile(user.Username, tempDirName, "app", language.Extension)
	if err != nil {
		rce.systemUsers.Release(user.Uid)
		internal.DeleteAll(user.Username)
		ch <- pool.CodeOutput{}
		return
	}

	// Write the code to the file.
	err = internal.WriteToFile(filename, data.Code)
	if err != nil {
		rce.systemUsers.Release(user.Uid)
		ch <- pool.CodeOutput{}
		return
	}

	executableFilename := internal.ExecutableFile(user.Username, tempDirName, "app")
	codeOutput := pool.CodeOutput{User: *user, TempDirName: tempDirName}

	// Compile the file when the language needs to be compiled.
	if language.Compiled {
		now := time.Now()
		compileOutput, compileError := rce.compileFile(filename, executableFilename, language)
		codeOutput.CompileOutput = pool.Output{
			Result: compileOutput,
			Error:  compileError,
			Time:   time.Since(now).Milliseconds(),
		}
	}

	// Execute the file when there is no error while compiling.
	if len(codeOutput.CompileOutput.Error) == 0 {
		now := time.Now()
		runOutput, runError := rce.executeFile(user.Username, filename, data.Stdin, executableFilename, language)
		codeOutput.RunOutput = pool.Output{
			Result: runOutput,
			Error:  runError,
			Time:   time.Since(now).Milliseconds(),
		}
	}

	// If the length of the test content is not empty, run the tests in the directory.
	if len(data.Tests) != 0 {
		now := time.Now()
		results := []pool.TestResult{}

		// Create a wait group to let the tests run concurrently and wait until all executed.
		var wg sync.WaitGroup
		wg.Add(len(data.Tests))
		for _, test := range data.Tests {
			go func(test pool.TestResult) {
				runOutput, runError := rce.executeFile(user.Username, filename, test.Stdin, executableFilename, language)
				if len(runError) != 0 {
					results = append(results, pool.TestResult{
						Name:     test.Name,
						Received: "",
						Actual:   test.Actual,
						Stdin:    test.Stdin,
						Passed:   false,
						RunError: runError,
					})
				} else {
					normalizedRunOutput := strings.TrimSuffix(runOutput, "\n")
					results = append(results, pool.TestResult{
						Name:     test.Name,
						Received: normalizedRunOutput,
						Actual:   test.Actual,
						Stdin:    test.Stdin,
						Passed:   test.Actual == normalizedRunOutput,
						RunError: "",
					})
				}
				wg.Done()
			}(test)
		}

		wg.Wait()

		codeOutput.TestOutput = pool.TestOutput{
			Results: results,
			Time:    time.Since(now).Milliseconds(),
		}
	}

	ch <- codeOutput

	if !data.BypassCache {
		rce.cache.Set(language.Name, data.Code, codeOutput)
	}

	rce.CleanUp(user, tempDirName)
}

// Dispatch dispatches a new job to the worker pool and returns the output of the
// submitted job.
func (rce *RceEngine) Dispatch(lang, code string, stdin []string, tests []pool.TestResult, bypassCache bool) (pool.CodeOutput, error) {
	dataChannel := make(chan pool.CodeOutput)
	rce.pool.SubmitJob(pool.WorkData{Lang: lang, Code: code, Stdin: stdin, Tests: tests, BypassCache: bypassCache}, rce.action, dataChannel)
	output := <-dataChannel
	return output, nil
}

// CleanUp cleans up the user's temporary directory, kills all running processes for this
// user, releases the user's uid and restores the user's directory.
func (rce *RceEngine) CleanUp(user *pool.User, tempDirName string) {
	internal.DeleteAll(user.Username)
	rce.cleanProcesses(user.Username)
	rce.restoreUserDir(user.Username)
	rce.systemUsers.Release(user.Uid)
}

// compileFile compiles the file and returns the output and possible error of the
// compilation. It uses the `compile.sh` script in the language directory to compile.
func (rce *RceEngine) compileFile(file, executableFile string, language Language) (string, string) {
	compileScript := fmt.Sprintf("/kira/languages/%s/compile.sh", strings.ToLower(language.Name))

	compile := exec.Command("/bin/bash", compileScript, file, executableFile)
	head := exec.Command("head", "--bytes", maxOutputBufferCapacity)

	errBuffer := bytes.Buffer{}
	compile.Stderr = &errBuffer

	head.Stdin, _ = compile.StdoutPipe()
	headOutput := bytes.Buffer{}
	head.Stdout = &headOutput

	_ = compile.Start()
	_ = head.Start()
	_ = compile.Wait()
	_ = head.Wait()

	result := ""

	if headOutput.Len() > 0 {
		result = headOutput.String()
	} else if headOutput.Len() == 0 && errBuffer.Len() == 0 {
		result = headOutput.String()
	}

	return result, errBuffer.String()
}

// executeFile executes the file and returns the output and possible error of the execution.
// It uses the `run.sh` script in the language directory to execute.
func (rce *RceEngine) executeFile(currentUser, file string, stdin []string, executableFile string, language Language) (string, string) {
	runScript := fmt.Sprintf("/kira/languages/%s/%s.sh", strings.ToLower(language.Name), "run")

	input := ""
	for _, in := range stdin {
		input += fmt.Sprintf("%q ", in)
	}
	input = strings.TrimSpace(input)

	run := exec.Command("/bin/bash", runScript, currentUser, fmt.Sprintf("%s %s", file, input), executableFile)
	head := exec.Command("head", "--bytes", maxOutputBufferCapacity)

	errBuffer := bytes.Buffer{}
	run.Stderr = &errBuffer

	head.Stdin, _ = run.StdoutPipe()
	headOutput := bytes.Buffer{}
	head.Stdout = &headOutput

	_ = run.Start()
	_ = head.Start()
	_ = run.Wait()
	_ = head.Wait()

	result := ""

	if headOutput.Len() > 0 {
		result = headOutput.String()
	} else if headOutput.Len() == 0 && errBuffer.Len() == 0 {
		result = headOutput.String()
	}

	return result, errBuffer.String()
}

// cleanProcesses cleans up all processes of the user by killing all processes of the user.
func (rce *RceEngine) cleanProcesses(currentUser string) error {
	return exec.Command("pkill", "-9", "-u", currentUser).Run()
}

// restoreUserDir restores the user directory if it was deleted.
// Creates the directory at `/tmp/<user>` by running the `mkdir` command as the user.
func (rce *RceEngine) restoreUserDir(currentUser string) {
	userDir := "/tmp/" + currentUser
	if _, err := os.ReadDir(userDir); err != nil {
		if os.IsNotExist(err) {
			_ = exec.Command("runuser", "-u", currentUser, "--", "mkdir", userDir).Run()
		}
	}
}
