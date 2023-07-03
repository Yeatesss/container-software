package core

import (
	"bytes"
	"path"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

var _ SoftwareFinder = &TomcatFindler{}

type TomcatFindler struct{}

func init() {
	if _, ok := Finders[WEB]; !ok {
		Finders[WEB] = make(map[string]SoftwareFinder)
	}
	Finders[WEB]["tomcat"] = NewTomcatFindler()
}
func NewTomcatFindler() *TomcatFindler {
	return &TomcatFindler{}
}

func (m TomcatFindler) Verify(c *Container, thisis func(*Process, SoftwareFinder)) bool {
	var hit bool

	_ = c.Processes.Range(func(_ int, ps *Process) (err error) {
		var cmdline *bytes.Buffer
		cmdline, err = ps.Cmdline()
		if err != nil {
			return
		}
		if cmdline.Len() > 0 && strings.Contains(cmdline.String(), "tomcat") {
			hit = true
			thisis(ps, &m)
			return nil
		}
		return
	})
	return hit
}

func (m TomcatFindler) GetSoftware(c *Container) ([]*Software, error) {
	var softwares []*Software
	_ = c.Processes.Range(func(_ int, ps *Process) (err error) {
		var software = &Software{
			Name:         "tomcat",
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

		software.Version, software.ConfigPath, err = getTomcatVersionAndConfig(c.EnvPath, ps)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, nil
}

func getTomcatVersionAndConfig(envPath string, ps *Process) (version string, config string, err error) {
	var (
		base              string
		findVersionStdout *bytes.Buffer
		runLine           []string
	)

	findVersionStdout, err = ps.Run(
		command.EnterProcessNsRun(ps.Pid(), []string{"find", "/", "-name", "version.sh"}),
	)
	if err != nil {
		return
	}
	findVersionRaw := findVersionStdout.Bytes()
	base, err = getTomcatBase(ps)
	if err != nil {
		return
	}

	if base != "" {
		runLine = append(runLine, path.Join(base, "./bin/version.sh"))
	}
	for len(findVersionRaw) > 0 {
		var val []byte
		val, findVersionRaw = command.NextField(findVersionRaw)
		runLine = append(runLine, string(val))
		findVersionRaw = command.NextLine(findVersionRaw)
	}
LOOP:
	for _, run := range runLine {
		var stdout []byte
		stdout, err := command.EnterProcessNsRun(ps.Pid(), []string{"env", envPath, "." + run}).CombinedOutput()
		if err != nil {
			continue
		}
		if strings.Contains(string(stdout), "Tomcat") {
			var (
				versionByte []byte
				val         []byte
			)
			for len(stdout) > 0 {
				val = command.ReadLine(stdout)
				if base == "" {
					base = path.Join(run, "../")
				}
				if bytes.Contains(val, []byte("Server number:")) {
					versionByte, stdout = command.ReadField(stdout, 3)
					version = string(versionByte)
					break LOOP
				}
				stdout = command.NextLine(stdout)

			}

		}
	}
	if base != "" {
		ImagineConf := path.Join(base, "conf/context.xml")
		stdout, err := command.EnterProcessNsRun(ps.Pid(), []string{"ls", "-l", ImagineConf}).CombinedOutput()
		if err == nil && len(stdout) > 0 {
			config = ImagineConf
		}
	}

	return

}

func getTomcatBase(ps *Process) (base string, err error) {
	var (
		cmdlineByte []byte
		stdout      *bytes.Buffer
	)
	stdout, err = ps.Cmdline()
	if err != nil {
		return
	}
	cmdlineByte = stdout.Bytes()
	for len(cmdlineByte) > 0 {
		var val []byte
		val, cmdlineByte = command.NextField(cmdlineByte)
		if strings.Contains(string(val), "catalina.base=") {
			exp := strings.Split(string(val), "=")
			if len(exp) == 2 && exp[1] != "" {
				base = exp[1]
				break
			}
		}
		if strings.Contains(string(val), "catalina.home=") {
			exp := strings.Split(string(val), "=")
			if len(exp) == 2 && exp[1] != "" {
				base = exp[1]
				break
			}
		}
	}
	return
}
