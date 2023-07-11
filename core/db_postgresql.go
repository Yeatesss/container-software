package core

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

var _ SoftwareFinder = &PostgresqlFindler{}

type PostgresqlFindler struct{}

func init() {
	if _, ok := Finders[DATABASE]; !ok {
		Finders[DATABASE] = make(map[string]SoftwareFinder)
	}
	Finders[DATABASE]["postgresql"] = PewpostgresqlFindler()
}
func PewpostgresqlFindler() *PostgresqlFindler {
	return &PostgresqlFindler{}
}

func (m PostgresqlFindler) Verify(c *Container, thisis func(*Process, SoftwareFinder)) bool {
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
		if len(exe) > 0 && strings.Contains(stdout.String(), "postgres") {
			hit = true
			thisis(ps, &m)
			return nil
		}
		return
	})
	return hit
}

func (m PostgresqlFindler) GetSoftware(c *Container) ([]*Software, error) {
	var softwares []*Software
	_ = c.Processes.Range(func(_ int, ps *Process) (err error) {
		var software = &Software{
			Name:         "postgresql",
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
		software.Version, err = getPostgresqlVersion(ps, exe)
		if err != nil {
			return err
		}
		software.ConfigPath, err = getPostgresqlConfig(ps)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, nil
}

func getPostgresqlVersion(ps process.Process, exe string) (string, error) {
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

func getPostgresqlConfig(ps process.Process) (string, error) {
	cmdline, err := ps.Cmdline()
	if err != nil {
		return "", err
	}
	if strings.Contains(cmdline.String(), "--config-file=") {
		cmdlineByte := cmdline.Bytes()
		for len(cmdlineByte) > 0 {
			var field []byte
			field, cmdlineByte = command.NextField(cmdlineByte)
			if bytes.Contains(field, []byte("--config-file=")) {
				exfield := strings.Split(string(field), "=")
				if len(exfield) > 1 {
					return exfield[1], nil
				}
			}

		}
		re := regexp.MustCompile(`--config-file=[\x20=]+(\S+)`)
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
		command.EnterProcessNsRun(ps.Pid(), []string{"find", "/", "-name", "pg_hba.conf"}),
	)
	if err != nil {
		return "", err
	}
	configRaw := stdout.Bytes()
	for len(configRaw) > 0 {
		var val []byte
		val, configRaw = command.ReadField(configRaw, 1)
		return string(val), nil
	}

	return "", nil
}