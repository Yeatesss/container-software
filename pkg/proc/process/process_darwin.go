package process

import (
	"bytes"
	"context"
	"os/exec"
)

var _ Process = &DarwinProcess{}

type DarwinProcess struct {
}

func (d DarwinProcess) Run(cmdFuncs ...func() (*exec.Cmd, context.CancelFunc)) (stdout *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (d DarwinProcess) EnterProcessNsRun(ctx context.Context, pid int64, cmdStrs []string, envs ...string) func() (*exec.Cmd, context.CancelFunc) {
	//TODO implement me
	panic("implement me")
}

func (d DarwinProcess) NewExecCommand(ctx context.Context, name string, arg ...string) func() (*exec.Cmd, context.CancelFunc) {
	//TODO implement me
	panic("implement me")
}

func (d DarwinProcess) Pid() int64 {
	//TODO implement me
	panic("implement me")
}

func (d DarwinProcess) ChildPids() []int64 {
	//TODO implement me
	panic("implement me")
}

func (d DarwinProcess) SetChildPids(int64s []int64) {
	//TODO implement me
	panic("implement me")
}

func (d DarwinProcess) Comm(ctx context.Context) (exe *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (d DarwinProcess) Cwd(ctx context.Context) (cwd *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (d DarwinProcess) Cmdline() (cmdline *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (d DarwinProcess) Exe(ctx context.Context) (exe *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func (d DarwinProcess) NsPids(ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}
func (d DarwinProcess) PidNamespace(_ context.Context) (exe *bytes.Buffer, err error) {
	//TODO implement me
	panic("implement me")
}

func NewProcess(_ int64, _ []int64) *DarwinProcess {
	return &DarwinProcess{}
}
