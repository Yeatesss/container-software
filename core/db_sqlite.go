package core

import (
	"bytes"
	"context"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"
	"github.com/Yeatesss/container-software/pkg/log"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

var _ SoftwareFinder = &SqliteFindler{}

const Sqlite SwName = "sqlite"

type SqliteFindler struct{}

func init() {
	if _, ok := Finders[DATABASE]; !ok {
		Finders[DATABASE] = make(map[SwName]SoftwareFinder)
	}
	Finders[DATABASE][Sqlite] = NewSqliteFindler()
}
func NewSqliteFindler() *SqliteFindler {
	return &SqliteFindler{}
}

func (m SqliteFindler) Verify(ctx context.Context, c *Container, thisis func(*Process, SoftwareFinder)) (bool, error) {
	var hit bool
	log.Logger.Debugf("Start verify sqlite:%s", c.Id)
	defer log.Logger.Debugf("Finish verify sqlite:%s", c.Id)
	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
		var exe string
		exe, err = process.GetProcessExe(ctx, ps.Process)
		if err != nil {
			return
		}
		stdout, err := ps.Run(
			ps.EnterProcessNsRun(ctx, ps.Pid(), []string{exe, "-help"}))
		if err != nil {
			return err
		}
		if len(exe) > 0 && strings.Contains(stdout.String(), "sqlite") {
			hit = true
			thisis(ps, &m)
			return nil
		}
		return
	})
	return hit, err
}

func (m SqliteFindler) GetSoftware(ctx context.Context, c *Container) ([]*Software, error) {
	var softwares []*Software
	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
		var software = &Software{
			Name:         "sqlite",
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
		software.Version, err = getSqliteVersion(ctx, ps, exe)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, err
}

func getSqliteVersion(ctx context.Context, ps process.Process, exe string) (string, error) {
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
	val, _ := command.ReadField(stdout.Bytes(), 1)
	return string(val), nil

}
