package docker

/*
下面代码处理dockerexec执行命令时返回err和out掺杂在一起的问题，返回时会带有ascii的问题
剥离开out/err
*/

// type ExecResult struct {
// 	StdOut   string
// 	StdErr   string
// 	ExitCode int
// }

// func Exec(ctx context.Context, containerID string, command []string) (types.IDResponse, error) {
// 	docker, err := client.NewEnvClient()
// 	if err != nil {
// 		return types.IDResponse{}, err
// 	}
// 	defer closer(docker)

// 	config := types.ExecConfig{
// 		AttachStderr: true,
// 		AttachStdout: true,
// 		Cmd:          command,
// 	}

// 	return docker.ContainerExecCreate(ctx, containerID, config)
// }

// func InspectExecResp(ctx context.Context, id string) (ExecResult, error) {
// 	var execResult ExecResult
// 	docker, err := client.NewEnvClient()
// 	if err != nil {
// 		return execResult, err
// 	}
// 	defer closer(docker)

// 	resp, err := docker.ContainerExecAttach(ctx, id, types.ExecConfig{})
// 	if err != nil {
// 		return execResult, err
// 	}
// 	defer resp.Close()

// 	// read the output
// 	var outBuf, errBuf bytes.Buffer
// 	outputDone := make(chan error)

// 	go func() {
// 		// StdCopy demultiplexes the stream into two buffers
// 		_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
// 		outputDone <- err
// 	}()

// 	select {
// 	case err := <-outputDone:
// 		if err != nil {
// 			return execResult, err
// 		}
// 		break

// 	case <-ctx.Done():
// 		return execResult, ctx.Err()
// 	}

// 	stdout, err := ioutil.ReadAll(&outBuf)
// 	if err != nil {
// 		return execResult, err
// 	}
// 	stderr, err := ioutil.ReadAll(&errBuf)
// 	if err != nil {
// 		return execResult, err
// 	}

// 	res, err := docker.ContainerExecInspect(ctx, id)
// 	if err != nil {
// 		return execResult, err
// 	}

// 	execResult.ExitCode = res.ExitCode
// 	execResult.StdOut = string(stdout)
// 	execResult.StdErr = string(stderr)
// 	return execResult, nil
// }
