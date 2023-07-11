package core

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

type SwType string

const (
	DATABASE SwType = "database"
	WEB      SwType = "web"
)

var Finders = make(map[SwType]map[string]SoftwareFinder)

// GetSoftware Get the application through the container process
func GetSoftware(c *Container) (softs []*Software, err error) {
	var finders = make(map[SoftwareFinder]*Container)
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()
	err = c.Processes.Range(func(_ int, process *Process) (rangerr error) {
		var (
			ctr *Container
			ok  bool
		)
		if ctr, ok = finders[process._finder]; !ok {
			ctr = &Container{
				Id:      c.Id,
				EnvPath: c.EnvPath,
			}
		}
		ctr.Processes = append(ctr.Processes, &Process{
			Process: process.Process,
		})
		finders[process._finder] = ctr
		return nil
	})
	for finderHandle, container := range finders {
		if finderHandle != nil {
			software, e := finderHandle.GetSoftware(container)
			if e != nil {
				return nil, e
			}
			softs = append(softs, software...)
		}

	}
	return
}

// SoftwareFinder Software finder,Verify is used to determine whether the current container has the software
// through GetSoftware to obtain specific software information
type SoftwareFinder interface {
	Verify(c *Container, thisis func(*Process, SoftwareFinder)) bool
	GetSoftware(c *Container) ([]*Software, error)
}
type Processes []*Process

func (l Processes) Range(f func(idx int, process *Process) error) (err error) {
	for idx, process := range l {
		err = f(idx, process)
		if err != nil {
			continue
		}
	}
	return
}
func (l Processes) Len() int {
	return len(l)
}

// Process Information about the processes in the container
type Process struct {
	process.Process
	_finder SoftwareFinder
}

// Container Container-related information
// EnvPath is the container and Path-related environment variables
type Container struct {
	Id        string
	EnvPath   string
	Processes Processes
}

// Software Information about software in containers
type Software struct {
	Name         string   `json:"name"`
	Type         SwType   `json:"type"`
	Version      string   `json:"version"`
	BindEndpoint []string `json:"bind_endpoint"`
	User         string   `json:"user"`
	BinaryPath   string   `json:"binary_path"`
	ConfigPath   string   `json:"config_path"`
}

func (p *Process) SetFinder(s SoftwareFinder) {
	p._finder = s
}

func GetRunUser(ps process.Process) (string, error) {
	var (
		stdout *bytes.Buffer
		nsPids []string
		err    error
	)
	nsPids, err = ps.NsPids()
	if err != nil {
		return "", err
	}
	if len(nsPids) > 0 {
		stdout, err = ps.Run(
			exec.Command("nsenter", "-t", strconv.FormatInt(ps.Pid(), 10), "--pid", "--uts", "--ipc", "--net", "--mount",
				"cat", fmt.Sprintf("/proc/%s/status", nsPids[len(nsPids)-1])),
			exec.Command("grep", "Uid"),
		)
		if err != nil {
			return "", err
		}
		uid, _ := command.ReadField(stdout.Bytes(), 2)
		if len(uid) > 0 {
			stdout, err = ps.Run(
				command.EnterProcessNsRun(ps.Pid(), []string{"getent", "passwd", string(uid)}),
			)
			if err != nil {
				return string(uid), nil
			}
			if stdout.Len() > 0 {
				return strings.Split(stdout.String(), ":")[0], nil
			}

		}

	}

	return "", nil
}

func GetEndpoint(ps process.Process) ([]string, error) {
	var (
		stdout    *bytes.Buffer
		err       error
		endpoints []string
	)

	stdout, err = ps.Run(
		exec.Command("nsenter", "-t", strconv.FormatInt(ps.Pid(), 10), "-n", "netstat", "-anp"),
		exec.Command("grep", `tcp\|udp`),
		exec.Command("grep", strconv.FormatInt(ps.Pid(), 10)),
		exec.Command("grep", `LISTEN`),
	)
	if err != nil {
		return []string{}, err
	}
	endpointRaw := stdout.Bytes()
	for {
		if len(endpointRaw) == 0 {
			break
		}
		var val []byte
		val, endpointRaw = command.ReadField(endpointRaw, 4)
		if len(val) > 0 {
			endpoints = append(endpoints, string(val))
		}
		endpointRaw = command.NextLine(endpointRaw)
	}

	return endpoints, nil

}
