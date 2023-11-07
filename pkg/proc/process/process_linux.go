package process

import (
	"bytes"
	"context"
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"

	"github.com/Yeatesss/container-software/pkg/command"
)

var _ Process = &LinuxProcess{}

type LinuxProcess struct {
	*command.CmdRuner
	comm     *bytes.Buffer
	cmd      *bytes.Buffer
	cwd      *bytes.Buffer
	exe      *bytes.Buffer
	nsPid    int64
	pid      int64
	childPid []int64
}

func (p *LinuxProcess) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(struct {
		Pid       int64
		ChildPids []int64
	}{Pid: p.pid, ChildPids: p.childPid})
}
func (p *LinuxProcess) NsPid() int64 {
	return p.nsPid
}

func (p *LinuxProcess) SetNsPid(nsPid int64) {
	p.nsPid = nsPid
	return
}

func NewProcess(pid int64, childPid []int64, options ...command.Option) *LinuxProcess {
	var cmdRunner = command.NewCmdRuner()
	if childPid == nil {
		childPid = []int64{}
	}
	for _, option := range options {
		cmdRunner = option(cmdRunner)
	}
	return &LinuxProcess{pid: pid, childPid: childPid, CmdRuner: cmdRunner}
}

func (p *LinuxProcess) Pid() int64 {
	return p.pid
}
func (p *LinuxProcess) ChildPids() []int64 {
	return p.childPid
}
func (p *LinuxProcess) SetChildPids(childPids []int64) {
	p.childPid = childPids
}
func (p *LinuxProcess) Comm(ctx context.Context) (exe *bytes.Buffer, err error) {
	if p.comm != nil {
		return p.comm, nil
	}
	p.comm, err = p.Run(
		p.NewExecCommand(ctx, "cat", fmt.Sprintf("/proc/%d/comm", p.pid)),
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
func (p *LinuxProcess) Cwd(ctx context.Context) (cwd *bytes.Buffer, err error) {
	if p.cwd != nil {
		return p.cwd, nil
	}
	cwd, err = p.Run(
		p.NewExecCommand(ctx, "ls", "-l", fmt.Sprintf("/proc/%d/cwd", p.pid)),
	)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(bytes.TrimSpace(cwd.Bytes())), nil
}
func (p *LinuxProcess) Exe(ctx context.Context) (exe *bytes.Buffer, err error) {
	if p.exe != nil {
		return p.exe, nil
	}
	exe, err = p.Run(
		p.NewExecCommand(ctx, "ls", "-l", fmt.Sprintf("/proc/%d/exe", p.pid)),
	)
	if err != nil {
		fmt.Println("exe error:", fmt.Sprintf("/proc/%d/exe", p.pid), err)
		return nil, err
	}
	return bytes.NewBuffer(bytes.TrimSpace(exe.Bytes())), nil
}
func (p *LinuxProcess) PidNamespace(_ context.Context) (exe *bytes.Buffer, err error) {
	var pidNamespace string
	nsBuf, err := p.Run(
		p.NewExecCommand(
			context.Background(), "ls", "-l", fmt.Sprintf("/proc/%d/ns", p.pid),
		),
	)
	if err != nil {
		return
	}
	nsBuf = command.Grep(nsBuf, "pid")
	if nsBuf.Len() > 0 {
		pidNs, _ := command.ReadField(nsBuf.Bytes(), 11)
		if len(pidNs) > 10 {
			pidNamespace = string(pidNs)[5 : len(pidNs)-1]
		}
	}

	return bytes.NewBuffer([]byte(pidNamespace)), nil
}
func (p *LinuxProcess) NsPids(ctx context.Context) ([]string, error) {
	var nsPids []string
	stdout, err := p.CmdRuner.Run(
		p.NewExecCommand(ctx, "grep", "NSpid", fmt.Sprintf("/proc/%d/status", p.pid)),
	)
	if err != nil {
		return nil, err
	}
	if stdout.Len() > 0 {
		nsPid := stdout.Bytes()
		for len(nsPid) > 0 {
			var val []byte
			val, nsPid = command.NextField(nsPid)
			if bytes.Contains(val, []byte("NSpid:")) {
				continue
			}
			nsPids = append(nsPids, string(val))
		}
	}
	if len(nsPids) == 0 && p.nsPid > 0 {
		nsPids = append(nsPids, fmt.Sprintf("%d", p.pid))
		nsPids = append(nsPids, fmt.Sprintf("%d", p.nsPid))
	}
	return nsPids, nil
}
