package sandbox

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

type Sandbox struct {
	ctx context.Context

	UUID   string
	Runner *runner

	cli         *client.Client
	ContainerID string

	SourceVolumePath string
	fileName         string
}

type Output struct {
	Error bool
	Body  string
}

var hostVolumePath = path.Join(os.Getenv("APP_CONTAINER_PATH"), "volume")

func NewSandbox(language string, code []byte) (*Sandbox, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	runner := &PythonRunner

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

	fileName := "code.py"
	filePath := path.Join(sourceVolumePath, fileName)
	err = ioutil.WriteFile(filePath, code, 0755)
	if err != nil {
		return nil, err
	}

	return &Sandbox{
		ctx:              ctx,
		UUID:             uid.String(),
		Runner:           runner,
		cli:              cli,
		SourceVolumePath: sourceVolumePath,
		fileName:         fileName,
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
			Image:      s.Runner.Image,
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
				NanoCPUs: s.Runner.MaxCPUs * 1000000000,
				Memory:   s.Runner.MaxMemory * 1024 * 1024,
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

	runOutput, err := s.Execute("", "")
	if err != nil {
		return output, err
	}

	output = append(output, runOutput)
	return output, nil
}

func (s *Sandbox) setupEnvironment() (*Output, error) {
	if len(s.Runner.BuildCmd) != 0 {
		return s.Execute("", "")
	}

	return &Output{}, nil
}

func (s *Sandbox) Execute(cmd, fileName string) (*Output, error) {
	if len(cmd) == 0 {
		if len(s.Runner.BuildCmd) > 0 {
			cmd = s.Runner.BuildCmd + " && " + s.Runner.RunCmd
		} else {
			cmd = s.Runner.RunCmd
		}
	}

	idResponse, err := s.cli.ContainerExecCreate(s.ctx, s.ContainerID, types.ExecConfig{
		Env:          s.Runner.Env,
		Cmd:          []string{"/bin/sh", "-c", cmd},
		Tty:          true,
		AttachStderr: true,
		AttachStdout: true,
		AttachStdin:  true,
		Detach:       true,
	})

	if err != nil {
		return &Output{
			Error: true,
			Body:  err.Error(),
		}, nil
	}

	hr, err := s.cli.ContainerExecAttach(s.ctx, idResponse.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		return &Output{
			Error: true,
			Body:  err.Error(),
		}, nil
	}
	defer hr.Close()

	scanner := bufio.NewScanner(hr.Conn)
	respBody := ""
	for scanner.Scan() {
		respBody += scanner.Text() + "\n"
	}

	return &Output{
		Error: false,
		Body:  respBody,
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
