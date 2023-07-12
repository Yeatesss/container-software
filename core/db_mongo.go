package core

import (
	"bytes"
	"context"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"
	"github.com/Yeatesss/container-software/pkg/log"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

var _ SoftwareFinder = &MongoFindler{}

type MongoFindler struct{}

const Mongo SwName = "mongo"

func init() {
	if _, ok := Finders[DATABASE]; !ok {
		Finders[DATABASE] = make(map[SwName]SoftwareFinder)
	}
	Finders[DATABASE][Mongo] = NewMongoFindler()
}
func NewMongoFindler() *MongoFindler {
	return &MongoFindler{}
}

func (m MongoFindler) Verify(ctx context.Context, c *Container, thisis func(*Process, SoftwareFinder)) (bool, error) {
	var hit bool

	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
		var exe string
		exe, err = process.GetProcessExe(ctx, ps.Process)
		if err != nil {
			return
		}
		stdout, err := ps.Run(
			ps.EnterProcessNsRun(ctx, ps.Pid(), []string{exe, "-h"}))
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
	return hit, err
}

func (m MongoFindler) GetSoftware(ctx context.Context, c *Container) ([]*Software, error) {
	var softwares []*Software
	log.Logger.Debugf("Start verify mongodb:%s", c.Id)
	defer log.Logger.Debugf("Finish verify mongodb:%s", c.Id)
	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
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
		software.Version, err = getMongoVersion(ctx, ps, exe)
		if err != nil {
			return err
		}
		software.ConfigPath, err = getMongoConfig(ctx, ps)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, err
}

func getMongoVersion(ctx context.Context, ps process.Process, exe string) (string, error) {
	var (
		stdout *bytes.Buffer
		err    error
	)
	stdout, err = ps.Run(
		ps.EnterProcessNsRun(ctx, ps.Pid(), []string{exe, "-version"}),
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

func getMongoConfig(ctx context.Context, ps process.Process) (string, error) {
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
		ps.EnterProcessNsRun(ctx, ps.Pid(), []string{"find", "/", "-name", "mongod.conf"}),
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
