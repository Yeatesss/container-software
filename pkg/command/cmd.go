package command

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type CmdRuner struct {
	cache map[string]string
}

func NewCmdRuner() CmdRuner {
	return CmdRuner{cache: make(map[string]string)}
}
func EnterProcessNsRun(pid int64, cmdStrs []string) *exec.Cmd {
	cmds := append([]string{"-t", strconv.FormatInt(pid, 10), "--pid", "--uts", "--ipc", "--net", "--mount"}, cmdStrs...)

	return exec.Command("nsenter", cmds...)
}

func (p CmdRuner) Run(cmdS ...*exec.Cmd) (stdout *bytes.Buffer, err error) {
	var cmdStr bytes.Buffer
	for _, cmd := range cmdS {
		cmdStr.WriteString(cmd.String())
	}
	defer func() {
		if stdout != nil {
			p.cache[cmdStr.String()] = stdout.String()
		}
	}()
	if v, ok := p.cache[cmdStr.String()]; ok {
		return bytes.NewBuffer([]byte(v)), nil
	}
	stdout, err = p.cmdPipeRun(cmdS...)
	if err = parseExitError(err); err != nil {
		return stdout, err
	}

	if strings.Contains(stdout.String(), fmt.Sprintf("%s: not found", cmdS[0].Path)) {
		return &bytes.Buffer{}, nil
	}
	stdout = bytes.NewBuffer(bytes.TrimSpace(stdout.Bytes()))
	return
}

func parseExitError(err error) error {
	e, ok := err.(*exec.ExitError)
	if ok && len(e.Stderr) == 0 {
		return nil
	}
	if ok {
		return errors.New(string(e.Stderr))
	}
	return err
}
func (p CmdRuner) cmdPipeRun(cmdS ...*exec.Cmd) (stdout *bytes.Buffer, err error) {
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
