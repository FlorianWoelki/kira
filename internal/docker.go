package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type ContainerPort struct {
	client *client.Client
}

func NewContainerPort(client *client.Client) *ContainerPort {
	return &ContainerPort{
		client: client,
	}
}

func (cp *ContainerPort) CreateContainer(ctx context.Context, image, uuid, sourceVolumePath string) (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	var networkMode container.NetworkMode
	createContainerResp, err := cp.client.ContainerCreate(ctx,
		&container.Config{
			Image:      image,
			Tty:        true,
			WorkingDir: "/runtime",
		},
		&container.HostConfig{
			NetworkMode: networkMode,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: workDir + "/" + sourceVolumePath,
					Target: "/runtime",
				},
			},
		}, nil, nil, uuid)

	if err != nil {
		return "", err
	}

	return createContainerResp.ID, nil
}

func (cp *ContainerPort) ExecuteCommand(ctx context.Context, containerId, cmd string, env []string) (types.HijackedResponse, error) {
	idResponse, err := cp.client.ContainerExecCreate(ctx, containerId, types.ExecConfig{
		Env:          env,
		Cmd:          []string{"/bin/sh", "-c", cmd},
		AttachStderr: true,
		AttachStdout: true,
	})

	if err != nil {
		return types.HijackedResponse{}, err
	}

	hr, err := cp.client.ContainerExecAttach(ctx, idResponse.ID, types.ExecStartCheck{})
	if err != nil {
		return types.HijackedResponse{}, err
	}

	return hr, nil
}

func (cp *ContainerPort) StartContainer(ctx context.Context, containerId string) error {
	return cp.client.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
}

func (cp *ContainerPort) RemoveContainer(ctx context.Context, containerId string) error {
	if err := cp.client.ContainerStop(ctx, containerId, nil); err != nil {
		return fmt.Errorf("Failed to stop container: %v", err)
	}

	if err := cp.client.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}); err != nil {
		return fmt.Errorf("Failed to remove container: %v", err)
	}

	return nil
}
