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

var _ SoftwareFinder = &JbossFindler{}

const Jboss SwName = "jboss"

type JbossFindler struct{}

func init() {
	if _, ok := Finders[WEB]; !ok {
		Finders[WEB] = make(map[SwName]SoftwareFinder)
	}
	Finders[WEB][Jboss] = NewJbossFindler()
}
func NewJbossFindler() *JbossFindler {
	return &JbossFindler{}
}

func (m JbossFindler) Verify(ctx context.Context, c *Container, thisis func(*Process, SoftwareFinder)) (bool, error) {
	var hit bool
	log.Logger.Debugf("Start verify jboss:%s", c.Id)
	defer log.Logger.Debugf("Finish verify jboss:%s", c.Id)
	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
		if ps._finder != nil {
			return nil
		}
		hit, err = m.SingleVerify(ctx, c.EnvPath, ps, thisis)
		return
	})
	return hit, err
}
func (m JbossFindler) SingleVerify(ctx context.Context, envPath string, ps *Process, thisis func(*Process, SoftwareFinder)) (hit bool, err error) {
	var (
		exe    string
		stdout *bytes.Buffer
	)
	exe, err = process.GetProcessExe(ctx, ps)
	if err != nil {
		return
	}
	stdout, err = ps.Run(
		ps.EnterProcessNsRun(ctx, ps.Pid(), []string{"env", envPath, "." + exe, "--version", "2>&1"}))
	if err != nil {
		return
	}
	if stdout.Len() > 0 && bytes.Contains(stdout.Bytes(), []byte("JBoss")) {
		hit = true
		thisis(ps, &m)
		return
	}
	return
}

func (m JbossFindler) GetSoftware(ctx context.Context, c *Container) ([]*Software, error) {
	var softwares []*Software
	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
		var software = &Software{
			Name:         "jboss",
			Type:         WEB,
			Version:      "",
			BindEndpoint: nil,
			User:         "",
			BinaryPath:   "",
			ConfigPath:   "",
		}
		var (
			exe string
		)
		exe, err = process.GetProcessExe(ctx, ps.Process)
		if err != nil {
			return
		}
		software.BinaryPath = exe
		for _, pid := range append([]int64{ps.Pid()}, ps.ChildPids()...) {
			endpoints, e := GetEndpoint(ctx, process.NewProcess(pid, nil))
			if e != nil {
				continue
			}
			software.BindEndpoint = append(software.BindEndpoint, endpoints...)
		}
		software.User, err = GetRunUser(ctx, ps, c.EnvPath)
		if err != nil {
			return err
		}

		software.Version, software.ConfigPath, err = getJbossVersionAndConfig(ctx, c.EnvPath, ps)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, err
}

func getJbossVersionAndConfig(ctx context.Context, envPath string, ps *Process) (version string, config string, err error) {
	var (
		exe string
	)
	exe, _ = process.GetProcessExe(ctx, ps)
	stdout, _ := ps.Run(
		ps.EnterProcessNsRun(ctx, ps.Pid(), []string{"env", envPath, "." + exe, "--version", "2>&1"}))
	if strings.Contains(stdout.String(), "WildFly") {
		re := regexp.MustCompile(`WildFly\s+Full\s+([\w\.]+)`)
		c := re.FindStringSubmatch(stdout.String())
		if len(c) > 1 {
			version = c[1]
		}
	}
	regexpExe, _ := regexp.Compile(`JBoss\s+(\d+\.\d+\.\d+)\.GA`)
	//regexp.find
	if version == "" {
		match := regexpExe.FindStringSubmatch(stdout.String())
		if len(match) == 2 {
			version = match[1]
		}
	}
	regexpExe, _ = regexp.Compile(`version\s+(\d+\.\d+\.\d+)`)
	//regexp.find
	if version == "" {
		match := regexpExe.FindStringSubmatch(stdout.String())
		if len(match) == 2 {
			version = match[1]
		}
	}
	stdout, err = ps.Run(
		ps.EnterProcessNsRun(ctx, ps.Pid(), []string{"find", "/", "-path", "/proc", "-prune", "-o", "-path", "/lib", "-prune", "-o", "-path", "/lib64", "-prune", "-o", "-name", "standalone.xml", "-print"}),
	)
	if err != nil {
		return
	}
	configRaw := stdout.Bytes()
	if len(configRaw) > 0 {
		var val []byte
		val, configRaw = command.ReadField(configRaw, 1)
		config = string(val)
	}

	return

}
