package core

import (
	"bytes"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

var _ SoftwareFinder = &RedisFindler{}

type RedisFindler struct{}

func init() {
	if _, ok := Finders[DATABASE]; !ok {
		Finders[DATABASE] = make(map[string]SoftwareFinder)
	}
	Finders[DATABASE]["redis"] = NewRedisFindler()
}
func NewRedisFindler() *RedisFindler {
	return &RedisFindler{}
}

func (m RedisFindler) Verify(c *Container, thisis func(*Process, SoftwareFinder)) bool {
	var hit bool

	_ = c.Processes.Range(func(_ int, ps *Process) (err error) {
		var exe string
		exe, err = process.GetProcessExe(ps.Process)
		if err != nil {
			return
		}
		stdout, err := ps.Run(
			command.EnterProcessNsRun(ps.Pid(), []string{exe, "-v"}))
		if err != nil {
			return err
		}
		if len(exe) > 0 && strings.Contains(stdout.String(), "Redis") {
			hit = true
			thisis(ps, &m)
			return nil
		}
		return
	})
	return hit
}

func (m RedisFindler) GetSoftware(c *Container) ([]*Software, error) {
	var softwares []*Software
	_ = c.Processes.Range(func(_ int, ps *Process) (err error) {
		var software = &Software{
			Name:         "redis",
			Type:         DATABASE,
			Version:      "",
			BindEndpoint: nil,
			User:         "",
			BinaryPath:   "",
			ConfigPath:   "",
		}
		var exe string
		exe, err = process.GetProcessExe(ps.Process)
		if err != nil {
			return
		}
		software.BinaryPath = exe
		software.BindEndpoint, err = GetEndpoint(ps)
		if err != nil {
			return err
		}
		software.User, err = GetRunUser(ps)
		if err != nil {
			return err
		}
		software.Version, err = getRedisVersion(ps, exe)
		if err != nil {
			return err
		}
		software.ConfigPath, err = getRedisConfig(ps)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, nil
}

func getRedisVersion(ps process.Process, exe string) (string, error) {
	var (
		stdout *bytes.Buffer
		err    error
	)
	stdout, err = ps.Run(
		command.EnterProcessNsRun(ps.Pid(), []string{exe, "-v"}),
	)
	if err != nil {
		return "", err
	}
	val, _ := command.ReadField(stdout.Bytes(), 3)
	if len(val) > 0 {
		exVersion := bytes.Split(val, []byte("="))
		if len(exVersion) == 2 {
			return string(exVersion[1]), nil

		}
	}
	return "", nil

}

func getRedisConfig(ps process.Process) (string, error) {
	var cmdlineByte []byte
	cmdline, err := ps.Cmdline()
	if err != nil {
		return "", err
	}
	comm, _ := ps.Comm()
	commIdx := bytes.Index(cmdline.Bytes(), comm.Bytes())
	if commIdx >= 0 {
		cmdlineByte = bytes.ReplaceAll(cmdline.Bytes()[commIdx:], comm.Bytes(), []byte{})
		for len(cmdlineByte) > 0 {
			var flag []byte
			flag, cmdlineByte = command.NextField(cmdlineByte)
			if process.IsPath(string(flag)) {
				return string(flag), nil
			}
		}
	}

	return "", nil
}
