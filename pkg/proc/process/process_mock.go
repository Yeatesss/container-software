package process

import (
	"bytes"
	"context"
	"os/exec"
)

var _ Process = &MockProcess{}

type MockProcess struct {
	pid      int64
	childPid []int64
	nsPids   []string
	cmdline  string // lighttpd -D -f /etc/lighttpd/lighttpd.conf
	cwd      string // lrwxrwxrwx 1 root root 0 Jun 20 17:28 /proc/2548908/cwd -> //
	comm     string // lighttpd
	exe      string // lrwxrwxrwx 1 root root 0 Jun 19 10:52 /proc/2548908/exe -> /usr/sbin/lighttpd
}

func (p *MockProcess) PidNamespace(_ context.Context) (exe *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (p *MockProcess) PidNs(_ context.Context) (exe *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (p *MockProcess) NsPid() int64 {
	//TODO implement me
	panic("implement me")
}

func (p *MockProcess) SetNsPid(nsPid int64) {
	//TODO implement me
	panic("implement me")
}

func (p *MockProcess) NewExecCommand(ctx context.Context, name string, arg ...string) func() (*exec.Cmd, context.CancelFunc) {
	//TODO implement me
	panic("implement me")
}

func (p *MockProcess) Run(cmdS ...func() (*exec.Cmd, context.CancelFunc)) (stdout *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (p *MockProcess) EnterProcessNsRun(ctx context.Context, pid int64, cmdStrs []string, envs ...string) func() (*exec.Cmd, context.CancelFunc) {
	//TODO implement me
	panic("implement me")
}

func (p *MockProcess) SetChildPids(int64s []int64) {
	//TODO implement me
	panic("implement me")
}

func (p *MockProcess) ChildPids() []int64 {
	return p.childPid
}

func (p *MockProcess) NsPids(ctx context.Context) ([]string, error) {
	return p.nsPids, nil
}

func (p *MockProcess) Pid() int64 {
	return p.pid
}

func (p *MockProcess) Comm(ctx context.Context) (exe *bytes.Buffer, err error) {
	return bytes.NewBuffer(bytes.TrimSpace([]byte(p.comm))), nil
}
func (p *MockProcess) Cmdline() (cmdline *bytes.Buffer, err error) {
	return bytes.NewBuffer(bytes.TrimSpace([]byte(p.cmdline))), nil

}
func (p *MockProcess) Cwd(ctx context.Context) (cwd *bytes.Buffer, err error) {
	return bytes.NewBuffer(bytes.TrimSpace([]byte(p.cwd))), nil
}

func (p *MockProcess) Exe(ctx context.Context) (exe *bytes.Buffer, err error) {
	return bytes.NewBuffer(bytes.TrimSpace([]byte(p.exe))), nil
}
