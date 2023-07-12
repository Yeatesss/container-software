package command

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Yeatesss/container-software/pkg/log"

	"github.com/pkg/errors"
)

var cmdTimeout = time.Duration(10)
var DefaultCmdRunner = CmdRuner{
	timeout: cmdTimeout * time.Second,
	cache:   make(map[string]string),
	lock:    sync.Mutex{},
}

type CmdRuner struct {
	timeout time.Duration
	cache   map[string]string
	lock    sync.Mutex
}

func (p *CmdRuner) NewExecCommand(ctx context.Context, name string, arg ...string) func() (*exec.Cmd, context.CancelFunc) {
	return func() (*exec.Cmd, context.CancelFunc) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, p.timeout)
		cmd := exec.CommandContext(ctx, name, arg...)
		return cmd, cancel
	}
}
func (p *CmdRuner) NewExecCommandWithEnv(ctx context.Context, name string, arg []string, envs ...string) func() (*exec.Cmd, context.CancelFunc) {
	return func() (*exec.Cmd, context.CancelFunc) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, p.timeout)
		cmd := exec.CommandContext(ctx, name, arg...)
		cmd.Env = envs
		return cmd, cancel
	}
}
func (p *CmdRuner) fromCache(cmd string) (string, bool) {
	p.lock.Lock()
	defer p.lock.Unlock()
	res, ok := p.cache[cmd]
	return res, ok
}

func (p *CmdRuner) setCache(cmd string, result string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.cache[cmd] = result
	return
}

type Option func(*CmdRuner) *CmdRuner

var WithTimeout = func(timeoutSecond int64) func(runer *CmdRuner) *CmdRuner {
	return func(runer *CmdRuner) *CmdRuner {
		runer.timeout = time.Duration(timeoutSecond) * time.Second
		return runer
	}
}

func NewCmdRuner(options ...Option) *CmdRuner {
	runner := &CmdRuner{
		timeout: cmdTimeout * time.Second,
		cache:   make(map[string]string),
		lock:    sync.Mutex{},
	}
	for _, option := range options {
		runner = option(runner)
	}
	return runner
}
func (p *CmdRuner) EnterProcessNsRun(ctx context.Context, pid int64, cmdStrs []string, envs ...string) func() (*exec.Cmd, context.CancelFunc) {
	cmds := append([]string{"-t", strconv.FormatInt(pid, 10), "--pid", "--uts", "--ipc", "--net", "--mount"}, cmdStrs...)
	return p.NewExecCommandWithEnv(ctx, "nsenter", cmds, envs...)
}
func Grep(stdout *bytes.Buffer, filters ...string) *bytes.Buffer {
	var res = &bytes.Buffer{} // 创建一个新的bytes.Buffer

	if stdout != nil {
		cmdByte := stdout.Bytes()
		for len(cmdByte) > 0 {
			line := ReadLine(cmdByte)
			for _, filter := range filters {
				cfilters := strings.Split(filter, "|")
				for _, cfilter := range cfilters {
					if bytes.Contains(line, []byte(cfilter)) {
						res.Write(line)     // 将符合预期的值写入res
						res.WriteByte('\n') // 添加换行符
					}
				}
			}
			cmdByte = NextLine(cmdByte)
		}
	}
	return res
}
func (p *CmdRuner) Run(cmdFuncs ...func() (*exec.Cmd, context.CancelFunc)) (stdout *bytes.Buffer, err error) {
	var cmdStr bytes.Buffer
	var cmdS []*exec.Cmd
	for _, f := range cmdFuncs {
		cmd, _ := f()
		cmdS = append(cmdS, cmd)
		cmdStr.WriteString(cmd.String())
	}
	defer func() {
		if stdout != nil {
			p.setCache(cmdStr.String(), stdout.String())
		}
	}()

	if v, ok := p.fromCache(cmdStr.String()); ok {
		return bytes.NewBuffer([]byte(v)), nil
	}

	stdout, err = p.cmdPipeRun(cmdS)

	if err != nil && stdout == nil {
		return nil, err
	}
	if err = parseExitError(err); err != nil {
		return stdout, err
	}

	if strings.Contains(stdout.String(), fmt.Sprintf("%s: not found", cmdS[0].Path)) {
		return &bytes.Buffer{}, nil
	}
	stdout = bytes.NewBuffer(bytes.TrimSpace(stdout.Bytes()))
	return
}

func parseExitError(err error) error {
	e, ok := err.(*exec.ExitError)
	if ok && len(e.Stderr) == 0 {
		return nil
	}
	if ok {
		return errors.New(string(e.Stderr))
	}
	return err
}

func (p *CmdRuner) cmdPipeRun(cmdS []*exec.Cmd) (out *bytes.Buffer, err error) {
	out = &bytes.Buffer{}
	var stdout io.ReadCloser
	var stderr io.ReadCloser
	defer func() {
		stdout.Close()
		stderr.Close()
	}()
	var bufByte []byte
	for idx := 0; idx < len(cmdS); idx++ {
		idx := idx
		if idx > 0 {
			cmdS[idx].Stdin = bytes.NewBuffer(bufByte)
		}
		stdout, err = cmdS[idx].StdoutPipe()
		if err != nil {
			return
		}
		stderr, err = cmdS[idx].StderrPipe()
		if err != nil {
			return
		}
		err = cmdS[idx].Start()
		if err != nil {
			return
		}
		out = &bytes.Buffer{}
		ret := make(chan error, 1)
		go CopyStd(ret, out, stdout)
		fuse := NewTimer(p.timeout)
		select {
		case e := <-ret:
			log.Logger.Debug("copy stdout success")
			if e != nil {
				log.Logger.Error("copy stdout fail:", e)
				return
			}
			if out.Len() == 0 {
				io.Copy(out, stderr)
			}
			bufByte = out.Bytes()
		case <-fuse.C:
			log.Logger.Debug("Get Stdout Timeout")
			return
		}

		if err = cmdS[idx].Wait(); err != nil {
			return nil, err
		}

	}
	return
}
func CopyStd(ret chan error, dst *bytes.Buffer, std io.ReadCloser) {
	defer close(ret)
	_, err := io.Copy(dst, std)
	ret <- err
}

func NewTimer(expire time.Duration) *time.Timer {
	return time.NewTimer(expire)

}
