package sandbox

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/florianwoelki/kira/internal"
	"github.com/google/uuid"
)

type SandboxFile struct {
	Code     []byte
	FileName string
}

type Sandbox struct {
	containerPort    *internal.ContainerPort
	ctx              context.Context
	UUID             string
	Language         *Language
	LastTimestamp    time.Time
	ContainerID      string
	SourceVolumePath string
	fileName         string
	tests            []SandboxFile
	forceQuit        bool
}

type Output struct {
	Error     bool
	ExecBody  string
	TestsBody string
}

var hostVolumePath = path.Join(os.Getenv("APP_CONTAINER_PATH"), "volume")

func NewSandbox(language string, mainCode []byte, files []SandboxFile, sandboxTests []SandboxFile) (*Sandbox, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	var lang *Language
	for _, l := range Languages {
		if language == l.Name {
			lang = &l
			break
		}
	}

	if lang == nil {
		return nil, fmt.Errorf("no language found with name %s", language)
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}

	cli.NegotiateAPIVersion(ctx)

	sourceVolumePath := path.Join(hostVolumePath, uid.String())
	err = os.MkdirAll(sourceVolumePath, 0755)
	if err != nil {
		return nil, err
	}

	fileName := "app" + lang.Ext
	filePath := path.Join(sourceVolumePath, fileName)
	err = ioutil.WriteFile(filePath, mainCode, 0755)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filePath := path.Join(sourceVolumePath, file.FileName)
		err = ioutil.WriteFile(filePath, file.Code, 0755)
		if err != nil {
			return nil, err
		}
	}

	testCommandAppendix := ""
	for i, test := range sandboxTests {
		testFilePath := path.Join(sourceVolumePath, test.FileName)
		err = ioutil.WriteFile(testFilePath, test.Code, 0755)
		if err != nil {
			return nil, err
		}

		testCommandAppendix += test.FileName
		if i != len(sandboxTests)-1 {
			testCommandAppendix += " "
		}
	}

	if len(sandboxTests) != 0 {
		lang.TestCommand = strings.Replace(lang.TestCommand, "{}", testCommandAppendix, 1)
		lang.TestCommand = strings.Replace(lang.TestCommand, "{}", fileName, 1)
	} else {
		lang.TestCommand = ""
	}

	return &Sandbox{
		ctx:              ctx,
		UUID:             uid.String(),
		Language:         lang,
		LastTimestamp:    time.Now(),
		SourceVolumePath: sourceVolumePath,
		fileName:         fileName,
		tests:            sandboxTests,
		containerPort:    internal.NewContainerPort(cli),
	}, nil
}

type RunOutput struct {
	SetupOutput *Output
	RunOutput   *Output
}

func (s *Sandbox) Run() (RunOutput, error) {
	id, err := s.containerPort.CreateContainer(s.ctx, s.Language.Image, s.UUID, s.SourceVolumePath)

	if err != nil {
		return RunOutput{}, err
	}

	s.ContainerID = id

	if err := s.containerPort.StartContainer(s.ctx, s.ContainerID); err != nil {
		return RunOutput{}, err
	}

	setupOutput, err := s.setupEnvironment()
	if err != nil {
		return RunOutput{}, err
	}

	output := RunOutput{SetupOutput: setupOutput}
	if setupOutput.Error {
		return output, nil
	}

	runOutput, err := s.Execute("", "", true)
	if err != nil {
		return output, err
	}

	output.RunOutput = runOutput
	return output, nil
}

func (s *Sandbox) setupEnvironment() (*Output, error) {
	if len(s.Language.BuildCmd) != 0 {
		return s.Execute("", "", false)
	}

	return &Output{}, nil
}

func (s *Sandbox) ExecCmdInSandbox(cmd string) (string, error) {
	hr, err := s.containerPort.ExecuteCommand(s.ctx, s.ContainerID, cmd, s.Language.Env)
	if err != nil {
		return "", err
	}
	defer hr.Close()

	scanner := bufio.NewScanner(hr.Conn)
	respBody := ""
	for scanner.Scan() {
		respBody += scanner.Text() + "\n"

		if s.forceQuit {
			break
		}
	}

	return respBody, nil
}

func (s *Sandbox) Execute(cmd, fileName string, executeTests bool) (*Output, error) {
	if len(cmd) == 0 {
		if len(s.Language.BuildCmd) > 0 {
			if len(fileName) > 0 && s.Language.Name == "java" && path.Ext(fileName) == s.Language.Ext {
				cmd = strings.ReplaceAll(s.Language.BuildCmd, s.Language.DefaultFileName, fileName) + " && " + strings.ReplaceAll(s.Language.RunCmd, "app", strings.ReplaceAll(fileName, s.Language.Ext, ""))
			} else {
				cmd = s.Language.BuildCmd + " && " + s.Language.RunCmd
			}
		} else {
			cmd = s.Language.RunCmd
		}
	}

	execOutput, err := s.ExecCmdInSandbox(cmd)
	if err != nil {
		return &Output{
			Error:     true,
			ExecBody:  err.Error(),
			TestsBody: "",
		}, err
	}

	testBody := ""
	if executeTests && len(s.Language.TestCommand) != 0 {
		testBody, err = s.ExecCmdInSandbox(s.Language.TestCommand)
		if err != nil {
			return &Output{
				Error:     true,
				ExecBody:  "",
				TestsBody: err.Error(),
			}, nil
		}
	}

	return &Output{
		Error:     false,
		ExecBody:  execOutput,
		TestsBody: testBody,
	}, nil
}

func (s *Sandbox) Clean() {
	if err := s.containerPort.RemoveContainer(s.ctx, s.ContainerID); err != nil {
		fmt.Print(err)
	}

	err := os.RemoveAll(s.SourceVolumePath)
	if err != nil {
		fmt.Printf("Failed to remove volume folder: %v", err)
	}
}
