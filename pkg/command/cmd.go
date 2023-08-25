package command

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Yeatesss/container-software/pkg/log"
	"golang.org/x/sync/singleflight"

	"github.com/pkg/errors"
)

var cmdTimeout = time.Duration(4)
var DefaultCmdRunner = CmdRuner{
	timeout: cmdTimeout * time.Second,
	cache:   make(map[string]string),
	lock:    sync.Mutex{},
}

type CmdRuner struct {
	timeout  time.Duration
	cache    map[string]string
	lock     sync.Mutex
	singleDo *singleflight.Group
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
func (p *CmdRuner) delCache(cmd string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.cache, cmd)
	return
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
		timeout:  cmdTimeout * time.Second,
		cache:    make(map[string]string),
		lock:     sync.Mutex{},
		singleDo: &singleflight.Group{},
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
			var match []byte
			for idx, filter := range filters {
			FILTER:
				if idx > 0 && len(match) == 0 {
					break
				}
				match = []byte{}
				cfilters := strings.Split(filter, "|")
				for _, cfilter := range cfilters {
					var next string
					if strings.Contains(cfilter, ">>") {
						excfilter := strings.Split(cfilter, ">>")
						cfilter = excfilter[0]
						next = excfilter[1]
					}
					if bytes.Contains(line, []byte(cfilter)) {
						match = line
						if next != "" {
							filter = next
							next = ""
							goto FILTER
						}
						break
					}
				}
			}
			if len(match) > 0 {
				res.Write(line)     // 将符合预期的值写入res
				res.WriteByte('\n') // 添加换行符
			}
			cmdByte = NextLine(cmdByte)
		}
	}
	return res
}
func (p *CmdRuner) CacheClear(cmdFuncs ...func() (*exec.Cmd, context.CancelFunc)) {
	var cmdStr bytes.Buffer
	var cmdS []*exec.Cmd
	for _, f := range cmdFuncs {
		cmd, _ := f()
		cmdS = append(cmdS, cmd)
		cmdStr.WriteString(cmd.String())
	}
	p.delCache(cmdStr.String())
}
func (p *CmdRuner) Run(cmdFuncs ...func() (*exec.Cmd, context.CancelFunc)) (stdout *bytes.Buffer, err error) {
	var cmdStr bytes.Buffer
	var cmdS []*exec.Cmd
	stdout = &bytes.Buffer{}
	for _, f := range cmdFuncs {
		cmd, _ := f()
		cmdS = append(cmdS, cmd)
		cmdStr.WriteString(cmd.String())
	}
	//fmt.Println(cmdStr.String())
	defer func() {
		if stdout != nil {
			p.setCache(cmdStr.String(), stdout.String())
		}
	}()

	if v, ok := p.fromCache(cmdStr.String()); ok {
		//fmt.Println("from cache", cmdStr.String())
		return bytes.NewBuffer([]byte(v)), nil
	}
	tmpStdout, doerr, _ := p.singleDo.Do(cmdStr.String(), func() (r interface{}, e error) {
		return p.cmdPipeRun(cmdS)
	})

	if doerr != nil && tmpStdout == nil {
		return nil, doerr
	}
	if doerr = parseExitError(doerr); doerr != nil {
		return stdout, doerr
	}
	stdout = tmpStdout.(*bytes.Buffer)
	if stdout == nil {
		return bytes.NewBuffer([]byte{}), nil
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
		go CopyStd(ret, out, stdout, stderr)
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
			//fmt.Println(cmdS[idx].Process.Kill())
			log.Logger.Debug("Get Stdout Timeout:", cmdS[idx].String())
			cmdS[idx].Process.Signal(os.Interrupt)

		}
		if err = cmdS[idx].Wait(); err != nil && out == nil {
			return nil, err
		}
	}
	err = nil
	return
}
func CopyStd(ret chan error, dst *bytes.Buffer, stds ...io.ReadCloser) {
	defer close(ret)
	for _, std := range stds {
		_, err := io.Copy(dst, std)
		if err != nil {
			ret <- err
			return
		}
		if dst != nil && dst.Len() > 0 {
			break
		}
	}
}

func NewTimer(expire time.Duration) *time.Timer {
	return time.NewTimer(expire)

}
