package internal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

const (
	maxOutputBufferCapacity = "65332"
)

type CodeOutput struct {
	User        User
	TempDirName string
	Result      string
	Error       string
}

type RceEngine struct {
	systemUsers *SystemUsers
}

func NewRceEngine() *RceEngine {
	return &RceEngine{
		systemUsers: NewSystemUser(50),
	}
}

func (rce *RceEngine) RunCode(lang, code string, retries int) (CodeOutput, error) {
	language, err := GetLanguageByName(lang)
	if err != nil {
		return CodeOutput{}, err
	}

	user, err := rce.systemUsers.Acquire()
	if err != nil {
		fmt.Println("error in acquire user", err)
	}

	tempDirName := uuid.New().String()

	err = CreateTempDir(user.username, tempDirName)
	if err != nil {
		if retries == 0 {
			return CodeOutput{}, nil
		}

		return rce.RunCode(lang, code, retries-1)
	}

	filename, err := CreateTempFile(user.username, tempDirName, language.Extension)
	if err != nil {
		if retries == 0 {
			return CodeOutput{}, nil
		}

		DeleteTempDir(user.username, tempDirName)
		return rce.RunCode(lang, code, retries-1)
	}

	err = WriteToFile(filename, code)
	if err != nil {
		return CodeOutput{}, err
	}

	output, errorString := rce.executeFile(user.username, filename, language)

	return CodeOutput{
		User:        *user,
		TempDirName: tempDirName,
		Result:      output,
		Error:       errorString,
	}, nil
}

func (rce *RceEngine) CleanUp(user User, tempDirName string) {
	DeleteTempDir(user.username, tempDirName)
	rce.cleanProcesses(user.username)
	rce.restoreUserDir(user.username)
	rce.systemUsers.Release(user.uid)
}

func (rce *RceEngine) executeFile(currentUser, file string, language Language) (string, string) {
	script := fmt.Sprintf("/kira/languages/%s/run.sh", strings.ToLower(language.Name))

	run := exec.Command("/bin/bash", script, currentUser, file)
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

	//fmt.Printf("user %s with error %s\n", currentUser, &errBuffer)

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
