package process

import (
	"bytes"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"

	"github.com/pkg/errors"
)

// GetExe Get the full path to the current process executable
func GetProcessExe(process Process) (string, error) {
	var (
		cmdline     *bytes.Buffer
		cmdlineByte []byte
		exePath     []byte
		commStr     string
	)
	comm, err := process.Comm()
	if err != nil {
		return "", err
	}
	commStr = strings.TrimSpace(comm.String())
	switch commStr {
	case "bash", "sh":
		cmdline, err = process.Cmdline()
		if err != nil {
			return "", err
		}
		cmdlineByte = cmdline.Bytes()
		for len(cmdlineByte) > 0 {
			exePath, cmdlineByte = command.ReadField(cmdlineByte, 1)
			if string(exePath) != commStr && IsPath(string(exePath)) {
				break
			}
		}
		if len(exePath) > 0 {
			if !filepath.IsAbs(string(exePath)) {
				var cwd *bytes.Buffer
				cwd, err = process.Cwd()
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
		exeBuf, err = process.Exe()
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
				command.EnterProcessNsRun(process.Pid(), []string{"find", "/", "-name", commStr}),
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
		return string(exe), nil

	}

	return "", nil
}

type Process interface {
	Run(cmdS ...*exec.Cmd) (stdout *bytes.Buffer, err error)
	Pid() int64
	ChildPids() []int64
	Comm() (comm *bytes.Buffer, err error)
	Cwd() (cwd *bytes.Buffer, err error)
	Cmdline() (cmdline *bytes.Buffer, err error)
	Exe() (exe *bytes.Buffer, err error)
	NsPids() ([]string, error)
}

func IsPath(path string) bool {
	path = strings.TrimSpace(path)
	if len(path) == 0 {
		return false
	}
	_, ok := specialCharacters[[]byte(path)[0]]
	return ok
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
