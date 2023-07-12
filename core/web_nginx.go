package core

import (
	"bytes"
	"context"
	"regexp"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"
	"github.com/Yeatesss/container-software/pkg/log"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

var _ SoftwareFinder = &NginxFindler{}

const Nginx SwName = "nginx"

type NginxFindler struct{}

func init() {
	if _, ok := Finders[WEB]; !ok {
		Finders[WEB] = make(map[SwName]SoftwareFinder)
	}
	Finders[WEB][Nginx] = NewNginxFindler()
}
func NewNginxFindler() *NginxFindler {
	return &NginxFindler{}
}

func (m NginxFindler) Verify(ctx context.Context, c *Container, thisis func(*Process, SoftwareFinder)) (bool, error) {
	var hit bool
	log.Logger.Debugf("Start verify nginx:%s", c.Id)
	defer log.Logger.Debugf("Finish verify nginx:%s", c.Id)
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
		if len(exe) > 0 && strings.Contains(stdout.String(), "nginx") {
			hit = true
			thisis(ps, &m)
			return nil
		}
		return
	})
	return hit, err
}

func (m NginxFindler) GetSoftware(ctx context.Context, c *Container) ([]*Software, error) {
	var softwares []*Software
	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
		var software = &Software{
			Name:         "nginx",
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
		software.User, err = GetRunUser(ctx, ps)
		if err != nil {
			return err
		}
		software.Version, err = getNginxVersion(ctx, ps, exe)
		if err != nil {
			return err
		}
		software.ConfigPath, err = getNginxConfig(ctx, ps, exe)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, err
}

func getNginxVersion(ctx context.Context, ps process.Process, exe string) (string, error) {
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
	val, _ := command.ReadField(stdout.Bytes(), 3)
	splitVal := bytes.Split(val, []byte("/"))
	if len(splitVal) > 1 {
		return string(splitVal[1]), nil
	}
	return "", nil

}

func getNginxConfig(ctx context.Context, ps process.Process, exe string) (string, error) {
	var (
		stdout *bytes.Buffer
		err    error
	)
	stdout, err = ps.Run(
		ps.EnterProcessNsRun(ctx, ps.Pid(), []string{exe, "-t"}),
	)
	if err != nil {
		return "", err
	}

	if strings.Contains(stdout.String(), ".conf") {
		re := regexp.MustCompile(`configuration\s+file\s+([\w/]+\.\w+)\s+test\s+is\s+successful`)
		c := re.FindStringSubmatch(stdout.String())
		if len(c) > 1 {
			return c[1], nil
		}
	}

	return "", nil
}
