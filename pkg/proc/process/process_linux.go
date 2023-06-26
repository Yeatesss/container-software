package process

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/Yeatesss/container-software/pkg/command"
)

var _ Process = &LinuxProcess{}

type LinuxProcess struct {
	comm *bytes.Buffer
	cmd  *bytes.Buffer
	cwd  *bytes.Buffer
	exe  *bytes.Buffer
	pid  int64
}

func NewProcess(pid int64) *LinuxProcess {
	return &LinuxProcess{pid: pid}
}
func (p *LinuxProcess) Pid() int64 {
	return p.pid
}
func (p *LinuxProcess) Comm() (exe *bytes.Buffer, err error) {
	if p.comm != nil {
		return p.comm, nil
	}
	p.comm, err = command.CmdRun(
		exec.Command("cat", fmt.Sprintf("/proc/%d/comm", p.pid)),
	)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(bytes.TrimSpace(p.comm.Bytes())), nil
}
func (p *LinuxProcess) Cmdline() (cmdline *bytes.Buffer, err error) {
	if p.cmd != nil {
		return p.cmd, nil
	}
	var tmpCmdline []byte
	tmpCmdline, err = os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", p.pid))
	if err != nil {
		return nil, err
	}
	p.cmd = bytes.NewBuffer(bytes.TrimSpace(tmpCmdline))
	return p.cmd, nil
}
func (p *LinuxProcess) Cwd() (cwd *bytes.Buffer, err error) {
	if p.cwd != nil {
		return p.cwd, nil
	}
	cwd, err = command.CmdRun(
		exec.Command("ls", "-l", fmt.Sprintf("/proc/%d/cwd", p.pid)),
	)
	return bytes.NewBuffer(bytes.TrimSpace(cwd.Bytes())), nil
}
func (p *LinuxProcess) Exe() (exe *bytes.Buffer, err error) {
	if p.exe != nil {
		return p.exe, nil
	}
	exe, err = command.CmdRun(
		exec.Command("ls", "-l", fmt.Sprintf("/proc/%d/exe", p.pid)),
	)
	return bytes.NewBuffer(bytes.TrimSpace(exe.Bytes())), nil
}
