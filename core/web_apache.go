package core

import (
	"bytes"
	"context"
	"path"
	"strings"

	"github.com/Yeatesss/container-software/pkg/log"

	"github.com/Yeatesss/container-software/pkg/command"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

var _ SoftwareFinder = &ApacheFindler{}

const Apache SwName = "apache"

type ApacheFindler struct{}

func init() {
	if _, ok := Finders[WEB]; !ok {
		Finders[WEB] = make(map[SwName]SoftwareFinder)
	}
	Finders[WEB][Apache] = NewApacheFindler()
}
func NewApacheFindler() *ApacheFindler {
	return &ApacheFindler{}
}

func (m ApacheFindler) Verify(ctx context.Context, c *Container, thisis func(*Process, SoftwareFinder)) (bool, error) {
	var hit bool
	log.Logger.Debugf("Start verify apache:%s", c.Id)
	defer log.Logger.Debugf("Finish verify apache:%s", c.Id)
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
		if len(exe) > 0 && strings.Contains(stdout.String(), "Apache") {
			hit = true
			thisis(ps, &m)
			return nil
		}
		return
	})
	return hit, err
}

func (m ApacheFindler) GetSoftware(ctx context.Context, c *Container) ([]*Software, error) {
	var softwares []*Software
	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
		var software = &Software{
			Name:         "apache",
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
		software.Version, err = getApacheVersion(ctx, ps, exe)
		if err != nil {
			return err
		}
		software.ConfigPath, err = getApacheConfig(ctx, ps, exe)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, err
}

func getApacheVersion(ctx context.Context, ps process.Process, exe string) (string, error) {
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

func getApacheConfig(ctx context.Context, ps process.Process, exe string) (string, error) {
	var (
		stdout              *bytes.Buffer
		err                 error
		cwd                 *bytes.Buffer
		getCfgPathByCmdline = func(buf *bytes.Buffer) string {
			if stdout.Len() > 0 && bytes.Contains(stdout.Bytes(), []byte("SERVER_CONFIG_FILE")) {
				cmdline := stdout.Bytes()
				basePath, _ := command.ReadField(cwd.Bytes(), 11)
				for len(cmdline) > 0 {
					var (
						flag   []byte
						exflag []string
					)
					flag, cmdline = command.NextField(cmdline)
					if bytes.Contains(flag, []byte("SERVER_CONFIG_FILE")) {
						exflag = strings.Split(string(flag), "=")
						if len(exflag) > 1 {
							return path.Join(string(basePath), strings.ReplaceAll(exflag[1], `"`, ""))
						}
					}
				}
			}
			return ""
		}
	)
	cwd, err = ps.Cwd(ctx)
	if err != nil {
		return "", err
	}
	stdout, err = ps.Cmdline()
	if err != nil {
		return "", err
	}
	cfg := getCfgPathByCmdline(stdout)
	if cfg != "" {
		return cfg, nil
	}
	stdout, err = ps.Run(
		ps.EnterProcessNsRun(ctx, ps.Pid(), []string{exe, "-V"}),
	)
	if err != nil {
		return "", err
	}
	cfg = getCfgPathByCmdline(stdout)
	if cfg != "" {
		return cfg, nil
	}
	stdout, err = ps.Run(
		ps.EnterProcessNsRun(ctx, ps.Pid(), []string{"find", "/", "-name", "httpd.conf"}),
	)
	if err != nil {
		return "", err
	}
	configRaw := stdout.Bytes()
	if len(configRaw) > 0 {
		cfgByte, _ := command.NextField(configRaw)
		return string(cfgByte), nil
	}
	return "", nil
}
