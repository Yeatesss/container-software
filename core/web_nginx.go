package core

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

var _ SoftwareFinder = &NginxFindler{}

type NginxFindler struct{}

func init() {
	if _, ok := Finders[WEB]; !ok {
		Finders[WEB] = make(map[string]SoftwareFinder)
	}
	Finders[WEB]["nginx"] = NewNginxFindler()
}
func NewNginxFindler() *NginxFindler {
	return &NginxFindler{}
}

func (m NginxFindler) Verify(c *Container, thisis func(*Process, SoftwareFinder)) bool {
	var hit bool

	_ = c.Processes.Range(func(_ int, ps *Process) (err error) {
		var exe string
		exe, err = process.GetProcessExe(ps.Process)
		if err != nil {
			return
		}
		stdout, err := command.CmdRun(
			command.EnterProcessNsRun(ps.Pid(), []string{exe, "-V"}))
		fmt.Println(stdout)
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
	return hit
}

func (m NginxFindler) GetSoftware(c *Container) ([]*Software, error) {
	var softwares []*Software
	_ = c.Processes.Range(func(_ int, ps *Process) (err error) {
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
		exe, err = process.GetProcessExe(ps.Process)
		if err != nil {
			return
		}
		software.BinaryPath = exe
		software.BindEndpoint, err = GetEndpoint(ps.Pid())
		if err != nil {
			return err
		}
		software.User, err = GetRunUser(ps.Pid())
		if err != nil {
			return err
		}
		software.Version, err = getNginxVersion(ps.Pid(), exe)
		if err != nil {
			return err
		}
		software.ConfigPath, err = getNginxConfig(ps.Pid(), exe)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, nil
}

func getNginxVersion(pid int64, exe string) (string, error) {
	var (
		stdout *bytes.Buffer
		err    error
	)
	stdout, err = command.CmdRun(
		command.EnterProcessNsRun(pid, []string{exe, "-V"}),
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

func getNginxConfig(pid int64, exe string) (string, error) {
	var (
		stdout *bytes.Buffer
		err    error
	)
	stdout, err = command.CmdRun(
		command.EnterProcessNsRun(pid, []string{exe, "-t"}),
	)
	if err != nil {
		return "", err
	}
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()
	if strings.Contains(stdout.String(), ".conf") {
		re := regexp.MustCompile(`configuration\s+file\s+([\w/]+\.\w+)\s+test\s+is\s+successful`)
		c := re.FindStringSubmatch(stdout.String())
		if len(c) > 1 {
			return c[1], nil
		}
	}

	return "", nil
}
