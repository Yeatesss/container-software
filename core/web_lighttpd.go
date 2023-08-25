package core

import (
	"bytes"
	"context"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"
	"github.com/Yeatesss/container-software/pkg/log"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

var _ SoftwareFinder = &LighttpdFindler{}

const Lighttpd SwName = "lighttpd"

type LighttpdFindler struct{}

func init() {
	if _, ok := Finders[WEB]; !ok {
		Finders[WEB] = make(map[SwName]SoftwareFinder)
	}
	Finders[WEB][Lighttpd] = NewLighttpdFindler()
}
func NewLighttpdFindler() *LighttpdFindler {
	return &LighttpdFindler{}
}

func (m LighttpdFindler) Verify(ctx context.Context, c *Container, thisis func(*Process, SoftwareFinder)) (bool, error) {
	var hit bool
	log.Logger.Debugf("Start verify lighttpd:%s", c.Id)
	defer log.Logger.Debugf("Finish verify lighttpd:%s", c.Id)
	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
		var exe string
		exe, err = process.GetProcessExe(ctx, ps.Process)
		if err != nil {
			return
		}
		stdout, err := ps.Run(
			ps.EnterProcessNsRun(ctx, ps.Pid(), []string{exe, "-V"}))
		if err != nil {
			return err
		}
		if len(exe) > 0 && strings.Contains(stdout.String(), "lighttpd") {
			hit = true
			thisis(ps, &m)
			return nil
		}
		return
	})
	return hit, err
}

func (m LighttpdFindler) GetSoftware(ctx context.Context, c *Container) ([]*Software, error) {
	var softwares []*Software
	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
		var software = &Software{
			Name:         "lighttpd",
			Type:         WEB,
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
		software.User, err = GetRunUser(ctx, ps, c.EnvPath)
		if err != nil {
			return err
		}
		software.Version, err = getLighttpdVersion(ctx, ps, exe)
		if err != nil {
			return err
		}
		software.ConfigPath, err = getLighttpdConfig(ctx, ps)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, err
}

func getLighttpdVersion(ctx context.Context, ps process.Process, exe string) (string, error) {
	var (
		stdout *bytes.Buffer
		err    error
	)
	stdout, err = ps.Run(
		ps.EnterProcessNsRun(ctx, ps.Pid(), []string{exe, "-V"}),
	)
	if err != nil {
		return "", err
	}
	val, _ := command.ReadField(stdout.Bytes(), 1)
	splitVal := bytes.Split(val, []byte("/"))
	if len(splitVal) > 1 {
		return string(splitVal[1]), nil
	}
	return "", nil

}

func getLighttpdConfig(ctx context.Context, ps process.Process) (string, error) {
	var configs []string
	cmdline, err := ps.Cmdline()
	if err != nil {
		return "", err
	}
	cmdlineByte := cmdline.Bytes()
	if strings.Contains(cmdline.String(), "-f") {
		var v []byte
		for len(cmdlineByte) > 0 {
			v, cmdlineByte = command.NextField(cmdlineByte)
			if string(v) == "-f" {
				v, cmdlineByte = command.NextField(cmdlineByte)
				return string(v), nil
			}
		}
	}
	var (
		stdout *bytes.Buffer
	)
	stdout, err = ps.Run(
		ps.EnterProcessNsRun(ctx, ps.Pid(), []string{"find", "/", "-name", "lighttpd.conf"}),
	)
	if err != nil {
		return "", err
	}
	configRaw := stdout.Bytes()
	for {
		if len(configRaw) == 0 {
			break
		}
		var val []byte
		val, configRaw = command.ReadField(configRaw, 1)
		if len(val) > 0 {
			configs = append(configs, string(val))
		}
		configRaw = command.NextLine(configRaw)
	}
	return strings.Join(configs, ","), nil
}
