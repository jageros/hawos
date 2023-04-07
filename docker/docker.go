package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
)

func Exec(ctx context.Context, container string, cmd []string, opts ...client.Opt) (string, error) {
	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return "", err
	}

	cid, err := cli.ContainerExecCreate(ctx, container, types.ExecConfig{
		User:         "root",
		Privileged:   true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
	})
	if err != nil {
		return "", err
	}

	rsp, err := cli.ContainerExecAttach(ctx, cid.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	defer rsp.Close()

	err = cli.ContainerExecStart(ctx, cid.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(rsp.Reader)
	return string(data), err
}

func Restart(ctx context.Context, containerID string, opts ...client.Opt) error {
	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return err
	}
	return cli.ContainerRestart(ctx, containerID, container.StopOptions{})
}
