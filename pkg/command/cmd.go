package command

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
)

func EnterProcessNsRun(pid int64, cmdStrs []string) *exec.Cmd {
	cmds := append([]string{"-t", strconv.FormatInt(pid, 10), "--pid", "--uts", "--ipc", "--net", "--mount"}, cmdStrs...)

	return exec.Command("nsenter", cmds...)
}

// type EnterNamespace func() string
//
//	var WithPidNs EnterNamespace = func() string {
//		return "--pid"
//	}
//
//	var WithUtcNs EnterNamespace = func() string {
//		return "--uts"
//	}
//
//	var WithIpcNs EnterNamespace = func() string {
//		return "--ipc"
//	}
//
//	var WithNetNs EnterNamespace = func() string {
//		return "--net"
//	}
//
//	var WithMountNs EnterNamespace = func() string {
//		return "--mount"
//	}
//
//	func MakeNsenterCmd(pid int64, enterNS ...EnterNamespace) (cmd *exec.Cmd,err error) {
//		if pid > 0 {
//			var nsflags = make([]string, 7, 7)
//			nsflags[0] = "-t"
//			nsflags[1] = strconv.FormatInt(pid, 10)
//			if len(enterNS) == 0 {
//				nsflags[2] = "--pid"
//				nsflags[3] = "--uts"
//				nsflags[4] = "--ipc"
//				nsflags[5] = "--net"
//				nsflags[6] = "--mount"
//			} else {
//				idx := 2
//				for _, ns := range enterNS {
//					nsflags[idx] = ns()
//					idx++
//				}
//			}
//			return exec.Command("nsenter", nsflags...)),nil
//
//		}
//		return nil,fmt.Errorf("pid is 0")
//	}
func CmdRun(cmdS ...*exec.Cmd) (stdout *bytes.Buffer, err error) {
	var stderr *bytes.Buffer
	stdout, stderr, err = CmdPipeRun(cmdS...)
	if err != nil {
		return
	}
	if stderr.Len() > 0 {
		err = fmt.Errorf(stderr.String())
		return
	}
	return
}
func CmdPipeRun(cmdS ...*exec.Cmd) (stdout, stderr *bytes.Buffer, err error) {
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
			stdout, stderr = SetCommandStd(cmd)
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

func SetCommandStd(cmd *exec.Cmd) (stdout, stderr *bytes.Buffer) {
	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return
}
