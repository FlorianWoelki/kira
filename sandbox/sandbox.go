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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

type SandboxTest struct {
	Code     []byte
	FileName string
}

type Sandbox struct {
	ctx              context.Context
	UUID             string
	Language         *Language
	LastTimestamp    time.Time
	cli              *client.Client
	ContainerID      string
	SourceVolumePath string
	fileName         string
	tests            []SandboxTest
}

type Output struct {
	Error     bool
	ExecBody  string
	TestsBody string
}

var hostVolumePath = path.Join(os.Getenv("APP_CONTAINER_PATH"), "volume")

func NewSandbox(language string, code []byte, sandboxTests []SandboxTest) (*Sandbox, error) {
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

	fileName := "code" + lang.Ext
	filePath := path.Join(sourceVolumePath, fileName)
	err = ioutil.WriteFile(filePath, code, 0755)
	if err != nil {
		return nil, err
	}

	for _, test := range sandboxTests {
		testFilePath := path.Join(sourceVolumePath, test.FileName)
		err = ioutil.WriteFile(testFilePath, test.Code, 0755)
		if err != nil {
			return nil, err
		}
	}

	return &Sandbox{
		ctx:              ctx,
		UUID:             uid.String(),
		Language:         lang,
		LastTimestamp:    time.Now(),
		cli:              cli,
		SourceVolumePath: sourceVolumePath,
		fileName:         fileName,
		tests:            sandboxTests,
	}, nil
}

func (s *Sandbox) Run() ([]*Output, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	output := make([]*Output, 0, 2)
	var networkMode container.NetworkMode
	createContainerResp, err := s.cli.ContainerCreate(s.ctx,
		&container.Config{
			Image:      s.Language.Image,
			Tty:        true,
			WorkingDir: "/runtime",
		},
		&container.HostConfig{
			NetworkMode: networkMode,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: workDir + "/" + s.SourceVolumePath,
					Target: "/runtime",
				},
			},
			Resources: container.Resources{
				NanoCPUs: s.Language.MaxCPUs * 1000000000,
				Memory:   s.Language.MaxMemory * 1024 * 1024,
			},
		}, nil, nil, s.UUID)

	if err != nil {
		return nil, err
	}

	s.ContainerID = createContainerResp.ID

	if err := s.cli.ContainerStart(s.ctx, s.ContainerID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	setupOutput, err := s.setupEnvironment()
	if err != nil {
		return output, err
	}

	output = append(output, setupOutput)
	if setupOutput.Error {
		return output, nil
	}

	runOutput, err := s.Execute("", "", true)
	if err != nil {
		return output, err
	}

	output = append(output, runOutput)
	return output, nil
}

func (s *Sandbox) setupEnvironment() (*Output, error) {
	if len(s.Language.BuildCmd) != 0 {
		return s.Execute("", "", false)
	}

	return &Output{}, nil
}

func (s *Sandbox) ExecCmdInSandbox(cmd string) (string, error) {
	idResponse, err := s.cli.ContainerExecCreate(s.ctx, s.ContainerID, types.ExecConfig{
		Env:          s.Language.Env,
		Cmd:          []string{"/bin/sh", "-c", cmd},
		Tty:          true,
		AttachStderr: true,
		AttachStdout: true,
		AttachStdin:  true,
		Detach:       true,
	})

	if err != nil {
		return "", err
	}

	hr, err := s.cli.ContainerExecAttach(s.ctx, idResponse.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		return "", err
	}
	defer hr.Close()

	scanner := bufio.NewScanner(hr.Conn)
	respBody := ""
	for scanner.Scan() {
		respBody += scanner.Text() + "\n"
	}

	return respBody, nil
}

func (s *Sandbox) Execute(cmd, fileName string, executeTests bool) (*Output, error) {
	if len(cmd) == 0 {
		if len(s.Language.BuildCmd) > 0 {
			if len(fileName) > 0 && s.Language.Name == "java" && path.Ext(fileName) == s.Language.Ext {
				cmd = strings.ReplaceAll(s.Language.BuildCmd, s.Language.DefaultFileName, fileName) + " && " + strings.ReplaceAll(s.Language.RunCmd, "code", strings.ReplaceAll(fileName, s.Language.Ext, ""))
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
	if err := s.cli.ContainerStop(s.ctx, s.ContainerID, nil); err != nil {
		fmt.Printf("Failed to stop container: %v", err)
	}

	if err := s.cli.ContainerRemove(s.ctx, s.ContainerID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}); err != nil {
		fmt.Printf("Failed to remove container: %v", err)
	}

	err := os.RemoveAll(s.SourceVolumePath)
	if err != nil {
		fmt.Printf("Failed to remove volume folder: %v", err)
	}
}
