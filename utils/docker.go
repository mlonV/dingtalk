package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// 传入用户，ip，密码来建立远程执行docker的链接 返回docker api
func NewDockerCli(user, addr, port string) (*client.Client, error) {
	helper, err := connhelper.GetConnectionHelper(fmt.Sprintf("ssh://%s@%s:%s", user, addr, port))
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: helper.Dialer,
		},
	}

	cli, err := client.NewClientWithOpts(
		client.WithHTTPClient(httpClient),
		client.WithHost(helper.Host),
		client.WithDialContext(helper.Dialer),
	)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// 获取容器所有信息
func GetContainerByDocker(cli *client.Client) ([]types.Container, error) {

	typesCon, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		// 加载所有的容器（退出的也包含在内）
		All: true,
	})
	if err != nil {
		return nil, err
	}
	return typesCon, nil
}

// nil
func ExecCmd(dockerCli *client.Client, cmd []string, conID string) (string, error) {
	dockerctx := context.Background()
	execConfig := types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		DetachKeys:   "ctrl-p,ctrl-q",
		Tty:          false,
		Cmd:          cmd,
		// Env: []string{
		// 	"FOO=bar",
		// 	"BAZ=quux",
		// },
	}

	id, err := dockerCli.ContainerExecCreate(dockerctx, conID, execConfig)
	if err != nil {
		return "", err
	}
	resp, err := dockerCli.ContainerExecAttach(dockerctx, id.ID, types.ExecStartCheck{Tty: false})
	if err != nil {
		return "", err
	}
	defer resp.Close()

	bufReader := bufio.NewReader(resp.Reader)
	buf := make([]byte, 1024)
	builder := strings.Builder{}

	for {
		n, err := bufReader.Read(buf)
		if err != nil || err == io.EOF || n == 0 {
			break
		}
		builder.Write(buf[:n])
	}
	return builder.String(), nil
}
