package process

import (
	"bytes"
	"context"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"

	"github.com/pkg/errors"
)

// GetExe Get the full path to the current process executable
func GetProcessExe(ctx context.Context, process Process) (string, error) {
	var (
		cmdline     *bytes.Buffer
		cmdlineByte []byte
		exePath     []byte
		commStr     string
	)
	comm, err := process.Comm(ctx)
	if err != nil {
		return "", err
	}
	commStr = strings.TrimSpace(comm.String())
	switch commStr {
	case "bash", "sh", "tini":
		var exePath []byte
		cmdline, err = process.Cmdline()
		if err != nil {
			return "", err
		}
		cmdlineByte = cmdline.Bytes()
		for len(cmdlineByte) > 0 {
			var tmpExePath []byte
			tmpExePath, cmdlineByte = command.ReadField(cmdlineByte, 1)
			if string(tmpExePath) != commStr && IsPath(string(tmpExePath)) {
				exePath = tmpExePath
			}
		}
		if len(exePath) > 0 {
			if !filepath.IsAbs(string(exePath)) {
				var cwd *bytes.Buffer
				cwd, err = process.Cwd(ctx)
				if err != nil {
					return "", errors.Wrap(err, "Failed to get process CWD")
				}
				if cwd.Len() == 0 {
					return "", nil
				}
				pick, _ := command.ReadField(cwd.Bytes(), 11)
				exePath = []byte(path.Join(string(pick), string(exePath)))
			}
			return string(exePath), nil
		}

	default:
		var exeBuf *bytes.Buffer
		exeBuf, err = process.Exe(ctx)
		if err != nil {
			return "", err
		}
		exe, _ := command.ReadField(exeBuf.Bytes(), 11)
		if strings.Contains(string(exe), "bash") || strings.Contains(string(exe), "sh") {
			//avoid /bin/bash comm -f xxx.sh
			cmdline, err = process.Cmdline()
			if err != nil {
				return "", err
			}
			cmdlineByte = cmdline.Bytes()
			for len(cmdlineByte) > 0 {
				exePath, cmdlineByte = command.ReadField(cmdlineByte, 1)
				if strings.Contains(string(exePath), commStr) {
					if path.IsAbs(string(exePath)) {
						return string(exePath), nil
					}
				}
			}
			var findStdOut *bytes.Buffer
			findStdOut, err = process.Run(
				process.EnterProcessNsRun(ctx, process.Pid(), []string{"find", "/", "-path", "/proc", "-prune", "-o", "-name", commStr, "-print"}),
			)
			if err != nil {
				return "", err
			}
			configRaw := findStdOut.Bytes()
			if len(configRaw) > 0 {
				var val []byte
				val, configRaw = command.ReadField(configRaw, 1)
				if bytes.Contains(val, comm.Bytes()) {
					return string(val), nil
				}

			}
			return "", nil

		}
		if strings.Contains(string(exe), "runsvinit") {
			return "blackexe", nil
		}
		return string(exe), nil

	}

	return "", nil
}

type Process interface {
	Run(cmdFuncs ...func() (*exec.Cmd, context.CancelFunc)) (stdout *bytes.Buffer, err error)
	CacheClear(cmdFuncs ...func() (*exec.Cmd, context.CancelFunc))
	EnterProcessNsRun(ctx context.Context, pid int64, cmdStrs []string, envs ...string) func() (*exec.Cmd, context.CancelFunc)
	NewExecCommand(ctx context.Context, name string, arg ...string) func() (*exec.Cmd, context.CancelFunc)
	NewExecCommandWithEnv(ctx context.Context, name string, arg []string, envs ...string) func() (*exec.Cmd, context.CancelFunc)
	NsPid() int64
	SetNsPid(nsPid int64)
	Pid() int64
	ChildPids() []int64
	SetChildPids([]int64)
	Comm(ctx context.Context) (exe *bytes.Buffer, err error)
	Cwd(ctx context.Context) (cwd *bytes.Buffer, err error)
	Cmdline() (cmdline *bytes.Buffer, err error)
	Exe(ctx context.Context) (exe *bytes.Buffer, err error)
	PidNamespace(_ context.Context) (exe *bytes.Buffer, err error)
	NsPids(ctx context.Context) ([]string, error)
}

func IsPath(pathStr string) bool {
	if path.IsAbs(pathStr) {
		return true
	}
	pathStr = strings.TrimSpace(pathStr)
	if len(pathStr) == 0 {
		return false
	}
	first := []byte(pathStr)[0]
	//a-zA-Z0-9.
	if (first >= 97 && first <= 122) || (first >= 65 && first <= 90) || (first >= 48 && first <= 57) || first == 46 {
		return true
	}
	return false
}

var specialCharacters = map[uint8]struct{}{
	126: {},
	33:  {},
	64:  {},
	35:  {},
	36:  {},
	37:  {},
	94:  {},
	38:  {},
	42:  {},
	40:  {},
	41:  {},
	95:  {},
	43:  {},
	123: {},
	125: {},
	124: {},
	58:  {},
	34:  {},
	60:  {},
	62:  {},
	63:  {},
	44:  {},
	46:  {},
	47:  {},
	59:  {},
	39:  {},
	91:  {},
	93:  {},
	92:  {},
}
