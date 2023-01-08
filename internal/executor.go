package internal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/florianwoelki/kira/internal/cache"
	"github.com/florianwoelki/kira/internal/pool"
	"github.com/google/uuid"
)

const (
	amountOfUsers           = 50
	maxOutputBufferCapacity = "65332"
)

type RceEngine struct {
	systemUsers *pool.SystemUsers
	pool        *pool.WorkerPool
	cache       *cache.Cache[pool.CodeOutput]
}

func NewRceEngine() *RceEngine {
	return &RceEngine{
		systemUsers: pool.NewSystemUser(amountOfUsers),
		pool:        pool.NewWorkerPool(amountOfUsers),
		cache:       cache.NewCache[pool.CodeOutput](),
	}
}

func (rce *RceEngine) action(data pool.WorkData, ch chan<- pool.CodeOutput) {
	language, err := GetLanguageByName(data.Lang)
	if err != nil {
		ch <- pool.CodeOutput{}
		return
	}

	var cacheOutput pool.CodeOutput
	if !data.BypassCache {
		cacheOutput, err = rce.cache.Get(language.Name, data.Code)

		if err == nil {
			ch <- cacheOutput
			return
		}
	}

	user, err := rce.systemUsers.Acquire()
	if err != nil {
		rce.systemUsers.Release(user.Uid)
		ch <- pool.CodeOutput{}
		return
	}

	tempDirName := uuid.New().String()

	err = CreateTempDir(user.Username, tempDirName)
	if err != nil {
		rce.systemUsers.Release(user.Uid)
		ch <- pool.CodeOutput{}
		return
	}

	filename, err := CreateTempFile(user.Username, tempDirName, "app", language.Extension)
	if err != nil {
		rce.systemUsers.Release(user.Uid)
		DeleteAll(user.Username)
		ch <- pool.CodeOutput{}
		return
	}

	err = WriteToFile(filename, data.Code)
	if err != nil {
		rce.systemUsers.Release(user.Uid)
		ch <- pool.CodeOutput{}
		return
	}

	testFilename := ""
	if len(data.Test) != 0 {
		testFilename, err = CreateTempFile(user.Username, tempDirName, "app_test", language.Extension)
		if err != nil {
			rce.systemUsers.Release(user.Uid)
			DeleteAll(user.Username)
			ch <- pool.CodeOutput{}
			return
		}

		err = WriteToFile(testFilename, data.Test)
		if err != nil {
			rce.systemUsers.Release(user.Uid)
			ch <- pool.CodeOutput{}
			return
		}
	}

	executableFilename := ExecutableFile(user.Username, tempDirName, "app")
	codeOutput := pool.CodeOutput{User: *user, TempDirName: tempDirName}

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
		runOutput, runError := rce.executeFile(user.Username, filename, executableFilename, language)
		codeOutput.RunOutput = pool.Output{
			Result: runOutput,
			Error:  runError,
			Time:   time.Since(now).Milliseconds(),
		}
	}

	// If the length of the test content is not empty, run the tests in the directory.
	if len(data.Test) != 0 {
		now := time.Now()
		executableFilename := ExecutableFile(user.Username, tempDirName, "code_test")
		testOutput, testError := rce.executeTestsForFile(user.Username, testFilename, executableFilename, language)
		codeOutput.TestOutput = pool.Output{
			Result: testOutput,
			Error:  testError,
			Time:   time.Since(now).Milliseconds(),
		}
	}

	ch <- codeOutput

	if !data.BypassCache {
		rce.cache.Set(language.Name, data.Code, codeOutput)
	}

	rce.CleanUp(user, tempDirName)
}

func (rce *RceEngine) Dispatch(lang, code string, test string, bypassCache bool) (pool.CodeOutput, error) {
	dataChannel := make(chan pool.CodeOutput)
	rce.pool.SubmitJob(pool.WorkData{Lang: lang, Code: code, Test: test, BypassCache: bypassCache}, rce.action, dataChannel)
	output := <-dataChannel
	return output, nil
}

func (rce *RceEngine) CleanUp(user *pool.User, tempDirName string) {
	DeleteAll(user.Username)
	rce.cleanProcesses(user.Username)
	rce.restoreUserDir(user.Username)
	rce.systemUsers.Release(user.Uid)
}

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

func (rce *RceEngine) executeTestsForFile(currentUser, file, executableFile string, language Language) (string, string) {
	return rce.execute(currentUser, file, "test", executableFile, language)
}

func (rce *RceEngine) executeFile(currentUser, file, executableFile string, language Language) (string, string) {
	return rce.execute(currentUser, file, "run", executableFile, language)
}

func (rce *RceEngine) execute(currentUser, file, scriptName, executableFile string, language Language) (string, string) {
	runScript := fmt.Sprintf("/kira/languages/%s/%s.sh", strings.ToLower(language.Name), scriptName)

	run := exec.Command("/bin/bash", runScript, currentUser, file, executableFile)
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

func (rce *RceEngine) cleanProcesses(currentUser string) error {
	return exec.Command("pkill", "-9", "-u", currentUser).Run()
}

func (rce *RceEngine) restoreUserDir(currentUser string) {
	userDir := "/tmp/" + currentUser
	if _, err := os.ReadDir(userDir); err != nil {
		if os.IsNotExist(err) {
			_ = exec.Command("runuser", "-u", currentUser, "--", "mkdir", userDir).Run()
		}
	}
}
