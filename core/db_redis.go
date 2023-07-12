package core

import (
	"bytes"
	"context"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"
	"github.com/Yeatesss/container-software/pkg/log"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

var _ SoftwareFinder = &RedisFindler{}

const Redis SwName = "redis"

type RedisFindler struct{}

func init() {
	if _, ok := Finders[DATABASE]; !ok {
		Finders[DATABASE] = make(map[SwName]SoftwareFinder)
	}
	Finders[DATABASE][Redis] = NewRedisFindler()
}
func NewRedisFindler() *RedisFindler {
	return &RedisFindler{}
}

func (m RedisFindler) Verify(ctx context.Context, c *Container, thisis func(*Process, SoftwareFinder)) (bool, error) {
	var hit bool
	log.Logger.Debugf("Start verify redis:%s", c.Id)
	defer log.Logger.Debugf("Finish verify redis:%s", c.Id)
	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
		var exe string
		exe, err = process.GetProcessExe(ctx, ps.Process)
		if err != nil {
			return
		}
		stdout, err := ps.Run(
			ps.EnterProcessNsRun(ctx, ps.Pid(), []string{exe, "-v"}))
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
	return hit, err
}

func (m RedisFindler) GetSoftware(ctx context.Context, c *Container) ([]*Software, error) {
	var softwares []*Software
	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
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
		exe, err = process.GetProcessExe(ctx, ps.Process)
		if err != nil {
			return
		}
		software.BinaryPath = exe
		software.BindEndpoint, err = GetEndpoint(ctx, ps)
		if err != nil {
			return err
		}
		software.User, err = GetRunUser(ctx, ps)
		if err != nil {
			return err
		}
		software.Version, err = getRedisVersion(ctx, ps, exe)
		if err != nil {
			return err
		}
		software.ConfigPath, err = getRedisConfig(ctx, ps)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, err
}

func getRedisVersion(ctx context.Context, ps process.Process, exe string) (string, error) {
	var (
		stdout *bytes.Buffer
		err    error
	)
	stdout, err = ps.Run(
		ps.EnterProcessNsRun(ctx, ps.Pid(), []string{exe, "-v"}),
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

func getRedisConfig(ctx context.Context, ps process.Process) (string, error) {
	var cmdlineByte []byte
	cmdline, err := ps.Cmdline()
	if err != nil {
		return "", err
	}
	comm, _ := ps.Comm(ctx)
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
