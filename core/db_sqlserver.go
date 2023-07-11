package core

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

var _ SoftwareFinder = &SqlServerFindler{}

type SqlServerFindler struct{}

func init() {
	if _, ok := Finders[DATABASE]; !ok {
		Finders[DATABASE] = make(map[string]SoftwareFinder)
	}
	Finders[DATABASE]["sqlserver"] = NewSqlServerFindler()
}
func NewSqlServerFindler() *SqlServerFindler {
	return &SqlServerFindler{}
}

func (m SqlServerFindler) Verify(c *Container, thisis func(*Process, SoftwareFinder)) bool {
	var hit bool

	_ = c.Processes.Range(func(_ int, ps *Process) (err error) {
		var exe string
		exe, err = process.GetProcessExe(ps.Process)
		if err != nil {
			return
		}
		cmd := command.EnterProcessNsRun(ps.Pid(), []string{exe, "-v"})
		cmd.Env = []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"}
		stdout, err := ps.Run(cmd)
		if err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				err = nil
			} else {
				return err
			}
		}
		if len(exe) > 0 && strings.Contains(stdout.String(), "SQL Server") {
			hit = true
			thisis(ps, &m)
			return nil
		}
		return
	})
	return hit
}

func (m SqlServerFindler) GetSoftware(c *Container) ([]*Software, error) {
	var softwares []*Software
	_ = c.Processes.Range(func(_ int, ps *Process) (err error) {
		var software = &Software{
			Name:         "sqlserver",
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
		software.Version, err = getSqlServerVersion(ps, exe)
		if err != nil {
			return err
		}
		software.ConfigPath = "-"
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, nil
}

func getSqlServerVersion(ps process.Process, exe string) (string, error) {
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
	val, _ := command.ReadField(stdout.Bytes(), 8)
	return string(val), nil

}
