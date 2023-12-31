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

var _ SoftwareFinder = &MysqlFindler{}

const Mysql SwName = "mysql"

type MysqlFindler struct{}

func init() {
	if _, ok := Finders[DATABASE]; !ok {
		Finders[DATABASE] = make(map[SwName]SoftwareFinder)
	}
	Finders[DATABASE][Mysql] = NewMysqlFindler()
}
func NewMysqlFindler() *MysqlFindler {
	return &MysqlFindler{}
}

func (m MysqlFindler) Verify(ctx context.Context, c *Container, thisis func(*Process, SoftwareFinder)) (bool, error) {
	var hit bool
	log.Logger.Debugf("Start verify mysql:%s", c.Id)
	defer log.Logger.Debugf("Finish verify mysql:%s", c.Id)
	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
		if ps._finder != nil {
			return nil
		}
		hit, err = m.SingleVerify(ctx, ps, thisis)
		return
	})
	return hit, err
}
func (m MysqlFindler) SingleVerify(ctx context.Context, ps *Process, thisis func(*Process, SoftwareFinder)) (hit bool, err error) {
	var (
		exe    string
		stdout *bytes.Buffer
	)
	exe, err = process.GetProcessExe(ctx, ps.Process)
	if err != nil {
		return
	}
	stdout, err = ps.Run(
		ps.EnterProcessNsRun(ctx, ps.Pid(), []string{exe, "-V"}))
	if err != nil {
		return
	}
	if len(exe) > 0 && strings.Contains(stdout.String(), "mysqld") {
		hit = true
		thisis(ps, &m)
		return
	}
	return
}
func (m MysqlFindler) GetSoftware(ctx context.Context, c *Container) ([]*Software, error) {
	var softwares []*Software
	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
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
		software.Version, err = getMysqlVersion(ctx, ps, exe)
		if err != nil {
			return err
		}
		software.ConfigPath, err = getMysqlConfig(ctx, ps)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, err
}

func getMysqlVersion(ctx context.Context, ps process.Process, exe string) (string, error) {
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
	return string(val), nil

}

func getMysqlConfig(ctx context.Context, ps process.Process) (string, error) {
	var configs []string
	cmdline, err := ps.Cmdline()
	if err != nil {
		return "", err
	}
	if strings.Contains(cmdline.String(), "--defaults-file") {
		re := regexp.MustCompile(`--defaults-file[\x20=]+([a-zA-Z0-9\_\-\/\.]+)`)
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
		ps.EnterProcessNsRun(ctx, ps.Pid(), []string{"find", "/", "-path", "/proc", "-prune", "-o", "-path", "/lib", "-prune", "-o", "-path", "/lib64", "-prune", "-o", "-name", "my.cnf", "-print"}),
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
