package process

import (
	"bytes"
)

var _ Process = &MockProcess{}

type MockProcess struct {
	pid     int64
	cmdline string // lighttpd -D -f /etc/lighttpd/lighttpd.conf
	cwd     string // lrwxrwxrwx 1 root root 0 Jun 20 17:28 /proc/2548908/cwd -> //
	comm    string // lighttpd
	exe     string // lrwxrwxrwx 1 root root 0 Jun 19 10:52 /proc/2548908/exe -> /usr/sbin/lighttpd
}

func (p *MockProcess) Pid() int64 {
	//TODO implement me
	panic("implement me")
}

func (p *MockProcess) Comm() (exe *bytes.Buffer, err error) {
	return bytes.NewBuffer(bytes.TrimSpace([]byte(p.comm))), nil
}
func (p *MockProcess) Cmdline() (cmdline *bytes.Buffer, err error) {
	return bytes.NewBuffer(bytes.TrimSpace([]byte(p.cmdline))), nil

}
func (p *MockProcess) Cwd() (cwd *bytes.Buffer, err error) {
	return bytes.NewBuffer(bytes.TrimSpace([]byte(p.cwd))), nil
}

func (p *MockProcess) Exe() (exe *bytes.Buffer, err error) {
	return bytes.NewBuffer(bytes.TrimSpace([]byte(p.exe))), nil
}