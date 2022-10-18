package internal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

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

func (rce *RceEngine) action(lang, code string, bypassCache bool, ch chan<- pool.CodeOutput) {
	language, err := GetLanguageByName(lang)
	if err != nil {
		ch <- pool.CodeOutput{}
		return
	}

	var cacheOutput pool.CodeOutput
	if !bypassCache {
		cacheOutput, err = rce.cache.Get(language.Name, code)

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

	filename, err := CreateTempFile(user.Username, tempDirName, language.Extension)
	if err != nil {
		rce.systemUsers.Release(user.Uid)
		DeleteAll(user.Username)
		ch <- pool.CodeOutput{}
		return
	}

	err = WriteToFile(filename, code)
	if err != nil {
		rce.systemUsers.Release(user.Uid)
		ch <- pool.CodeOutput{}
		return
	}

	executableFilename := ExecutableFile(user.Username, tempDirName)

	var compileOutput, compileErrorString string
	var runOutput, runErrorString string

	if language.Compiled {
		compileOutput, compileErrorString = rce.compileFile(filename, executableFilename, language)
	}

	if len(compileErrorString) == 0 {
		runOutput, runErrorString = rce.executeFile(user.Username, filename, executableFilename, language)
	}

	codeOutput := pool.CodeOutput{
		User:          *user,
		TempDirName:   tempDirName,
		CompileResult: compileOutput,
		CompileError:  compileErrorString,
		RunResult:     runOutput,
		RunError:      runErrorString,
	}

	ch <- codeOutput

	if !bypassCache {
		rce.cache.Set(language.Name, code, codeOutput)
	}

	rce.CleanUp(user, tempDirName)
}

func (rce *RceEngine) Dispatch(lang, code string, bypassCache bool) (pool.CodeOutput, error) {
	dataChannel := make(chan pool.CodeOutput)
	rce.pool.SubmitJob(lang, code, bypassCache, rce.action, dataChannel)
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

func (rce *RceEngine) executeFile(currentUser, file, executableFile string, language Language) (string, string) {
	runScript := fmt.Sprintf("/kira/languages/%s/run.sh", strings.ToLower(language.Name))

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
