package core

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

var _ SoftwareFinder = &MysqlFindler{}

type MysqlFindler struct{}

func init() {
	if _, ok := Finders[DATABASE]; !ok {
		Finders[DATABASE] = make(map[string]SoftwareFinder)
	}
	Finders[DATABASE]["mysql"] = NewMysqlFindler()
}
func NewMysqlFindler() *MysqlFindler {
	return &MysqlFindler{}
}

func (m MysqlFindler) Verify(c *Container, thisis func(*Process, SoftwareFinder)) bool {
	var hit bool

	_ = c.Processes.Range(func(_ int, ps *Process) (err error) {
		var exe string
		exe, err = process.GetProcessExe(ps.Process)
		if err != nil {
			return
		}
		stdout, err := ps.Run(
			command.EnterProcessNsRun(ps.Pid(), []string{exe, "-V"}))
		if err != nil {
			return err
		}
		if len(exe) > 0 && strings.Contains(stdout.String(), "mysqld") {
			hit = true
			thisis(ps, &m)
			return nil
		}
		return
	})
	return hit
}

func (m MysqlFindler) GetSoftware(c *Container) ([]*Software, error) {
	var softwares []*Software
	_ = c.Processes.Range(func(_ int, ps *Process) (err error) {
		var software = &Software{
			Name:         "mysql",
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
		software.Version, err = getMysqlVersion(ps, exe)
		if err != nil {
			return err
		}
		software.ConfigPath, err = getMysqlConfig(ps)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, nil
}

func getMysqlVersion(ps process.Process, exe string) (string, error) {
	var (
		stdout *bytes.Buffer
		err    error
	)
	stdout, err = ps.Run(
		command.EnterProcessNsRun(ps.Pid(), []string{exe, "-V"}),
	)
	if err != nil {
		return "", err
	}
	val, _ := command.ReadField(stdout.Bytes(), 3)
	return string(val), nil

}

func getMysqlConfig(ps process.Process) (string, error) {
	var configs []string
	cmdline, err := ps.Cmdline()
	if err != nil {
		return "", err
	}
	if strings.Contains(cmdline.String(), "--defaults-file") {
		re := regexp.MustCompile(`--defaults-file[\x20=]+(\S+)`)
		match := re.FindStringSubmatch(cmdline.String())

		if len(match) > 1 {
			return strings.TrimSpace(match[1]), nil
		}
		return "", nil
	}
	var (
		stdout *bytes.Buffer
	)
	stdout, err = ps.Run(
		command.EnterProcessNsRun(ps.Pid(), []string{"find", "/", "-name", "my.cnf"}),
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
