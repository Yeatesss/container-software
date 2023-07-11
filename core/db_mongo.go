package core

import (
	"bytes"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

var _ SoftwareFinder = &MongoFindler{}

type MongoFindler struct{}

func init() {
	if _, ok := Finders[DATABASE]; !ok {
		Finders[DATABASE] = make(map[string]SoftwareFinder)
	}
	Finders[DATABASE]["mongo"] = NewMongoFindler()
}
func NewMongoFindler() *MongoFindler {
	return &MongoFindler{}
}

func (m MongoFindler) Verify(c *Container, thisis func(*Process, SoftwareFinder)) bool {
	var hit bool

	_ = c.Processes.Range(func(_ int, ps *Process) (err error) {
		var exe string
		exe, err = process.GetProcessExe(ps.Process)
		if err != nil {
			return
		}
		stdout, err := ps.Run(
			command.EnterProcessNsRun(ps.Pid(), []string{exe, "-h"}))
		if err != nil {
			return err
		}
		if len(exe) > 0 && strings.Contains(stdout.String(), "mongodb") {
			hit = true
			thisis(ps, &m)
			return nil
		}
		return
	})
	return hit
}

func (m MongoFindler) GetSoftware(c *Container) ([]*Software, error) {
	var softwares []*Software
	_ = c.Processes.Range(func(_ int, ps *Process) (err error) {
		var software = &Software{
			Name:         "mongo",
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
		software.Version, err = getMongoVersion(ps, exe)
		if err != nil {
			return err
		}
		software.ConfigPath, err = getMongoConfig(ps)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, nil
}

func getMongoVersion(ps process.Process, exe string) (string, error) {
	var (
		stdout *bytes.Buffer
		err    error
	)
	stdout, err = ps.Run(
		command.EnterProcessNsRun(ps.Pid(), []string{exe, "-version"}),
	)
	if err != nil {
		return "", err
	}
	versionStdout := stdout.Bytes()
	for len(versionStdout) > 0 {
		line := command.ReadLine(versionStdout)
		if bytes.Contains(line, []byte("db version")) {
			return strings.TrimSpace(strings.Replace(string(line), "db version", "", -1)), nil
		}
	}
	return "", nil

}

func getMongoConfig(ps process.Process) (string, error) {
	cmdline, err := ps.Cmdline()
	if err != nil {
		return "", err
	}
	cmdlineByte := cmdline.Bytes()
	if strings.Contains(cmdline.String(), "-f") || strings.Contains(cmdline.String(), "--config") {
		for len(cmdlineByte) > 0 {
			var flag []byte
			flag, cmdlineByte = command.NextField(cmdlineByte)
			if string(flag) == "-f" || string(flag) == "--config" {
				flag, _ = command.NextField(cmdlineByte)
				return string(flag), nil
			}
		}
	}
	var (
		stdout *bytes.Buffer
	)
	stdout, err = ps.Run(
		command.EnterProcessNsRun(ps.Pid(), []string{"find", "/", "-name", "mongod.conf"}),
	)
	if err != nil {
		return "", err
	}
	configRaw := stdout.Bytes()
	for len(configRaw) > 0 {
		var val []byte
		val, configRaw = command.ReadField(configRaw, 1)
		if len(val) > 0 {
			return string(val), nil
		}
		configRaw = command.NextLine(configRaw)
	}
	return "-", nil
}
