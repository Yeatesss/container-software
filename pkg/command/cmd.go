package command

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

func EnterProcessNsRun(pid int64, cmdStrs []string) *exec.Cmd {
	cmds := append([]string{"-t", strconv.FormatInt(pid, 10), "--pid", "--uts", "--ipc", "--net", "--mount"}, cmdStrs...)

	return exec.Command("nsenter", cmds...)
}

func CmdRun(cmdS ...*exec.Cmd) (stdout *bytes.Buffer, err error) {
	stdout, err = CmdPipeRun(cmdS...)
	if err != nil {
		return nil, err
	}
	if strings.Contains(stdout.String(), fmt.Sprintf("%s: not found", cmdS[0].Path)) {
		return &bytes.Buffer{}, nil
	}
	return
}
func CmdPipeRun(cmdS ...*exec.Cmd) (stdout *bytes.Buffer, err error) {
	var out io.ReadCloser
	var in io.WriteCloser
	for i, cmd := range cmdS {
		if i > 0 {
			// 上一个指令的输出
			out, err = cmdS[i-1].StdoutPipe()
			if err != nil {
				return
			}
			// 当前指令的输入
			in, err = cmd.StdinPipe()
			if err != nil {
				return
			}
			// 阻塞指令进行，直到 Close()
			go cmdPipe(cmdS[i-1], out, in)
		}
		if i == len(cmdS)-1 {
			stdout = SetCommandStd(cmd)
		}
	}
	for _, cmd := range cmdS {
		err = cmd.Run()
		if err != nil {
			return
		}
	}
	return
}

func cmdPipe(cmd *exec.Cmd, r io.ReadCloser, w io.WriteCloser) {
	defer func() {
		_ = r.Close()
		_ = w.Close()
	}()
	_, err := io.Copy(w, r)
	if err != nil {
		fmt.Println("error: ", err)
	}
	return
}

func SetCommandStd(cmd *exec.Cmd) (stdout *bytes.Buffer) {
	stdout = &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stdout
	return
}
