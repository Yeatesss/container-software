package core

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/Yeatesss/container-software/pkg/command"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

type SwType string
type SwName string

const (
	DATABASE SwType = "database"
	WEB      SwType = "web"
)

var Finders = make(map[SwType]map[SwName]SoftwareFinder)

// GetSoftware Get the application through the container process
func GetSoftware(ctx context.Context, c *Container) (softs []*Software, err error) {
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
			software, e := finderHandle.GetSoftware(ctx, container)
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
	Verify(ctx context.Context, c *Container, thisis func(*Process, SoftwareFinder)) (bool, error)
	GetSoftware(ctx context.Context, c *Container) ([]*Software, error)
}
type Processes []*Process

func (l Processes) Range(f func(idx int, process *Process) error) (hasErr error) {
	var (
		skipChilds = make(map[int64]bool, l.Len())
	)
	sort.SliceStable(l, func(i, j int) bool {
		return l[i].Pid() < l[j].Pid()
	})
	for idx, ps := range l {
		if _, ok := skipChilds[ps.Pid()]; ok {
			continue
		}
		err := f(idx, ps)
		if err != nil {
			hasErr = err
			continue
		}

		//Ignore child process fetching if complete data can be fetched from parent process
		if len(ps.ChildPids()) > 0 {
			for _, childPid := range ps.ChildPids() {
				skipChilds[childPid] = true
			}
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

func GetRunUser(ctx context.Context, ps process.Process) (string, error) {
	var (
		stdout *bytes.Buffer
		nsPids []string
		ids    []string
		err    error
	)
	nsPids, err = ps.NsPids(ctx)
	if err != nil {
		return "", err
	}
	if len(nsPids) == 0 {
		nsPids = append(nsPids, strconv.FormatInt(ps.Pid(), 10))
		var idx = 2
		for _, uidType := range []string{"Uid", "Gid"} {
			idx++
			stdout, err = ps.Run(
				ps.NewExecCommand(ctx, "nsenter", "-t", strconv.FormatInt(ps.Pid(), 10), "--pid", "--uts", "--ipc", "--net",
					"cat", fmt.Sprintf("/proc/%s/status", nsPids[len(nsPids)-1])),
			)
			if err != nil {
				return "", err
			}
			stdout = command.Grep(stdout, uidType)
			id, _ := command.ReadField(stdout.Bytes(), 2)
			if len(id) > 0 {
				stdout, err = ps.Run(
					ps.EnterProcessNsRun(ctx, ps.Pid(), []string{"getent", "passwd"}),
					ps.NewExecCommand(ctx, "awk", "-F:", fmt.Sprintf(`$%d==%s{print}`, idx, string(id))),
				)
				if err != nil {
					ids = append(ids, string(id))
					continue
				}
				if stdout.Len() > 0 {
					id = []byte(strings.Split(stdout.String(), ":")[0])
				}

			}
			ids = append(ids, string(id))
		}
	} else {
		var idx = 2
		for _, uidType := range []string{"Uid", "Gid"} {
			idx++
			stdout, err = ps.Run(
				ps.NewExecCommand(ctx, "nsenter", "-t", strconv.FormatInt(ps.Pid(), 10), "--pid", "--uts", "--ipc", "--net", "--mount",
					"cat", fmt.Sprintf("/proc/%s/status", nsPids[len(nsPids)-1])),
			)
			if err != nil {
				return "", err
			}
			stdout = command.Grep(stdout, uidType)
			id, _ := command.ReadField(stdout.Bytes(), 2)
			if len(id) > 0 {
				stdout, err = ps.Run(
					ps.EnterProcessNsRun(ctx, ps.Pid(), []string{"getent", "passwd"}),
					ps.NewExecCommand(ctx, "awk", "-F:", fmt.Sprintf(`$%d==%s{print}`, idx, string(id))),
				)
				if err != nil {
					ids = append(ids, string(id))
					continue
				}
				if stdout.Len() > 0 {
					id = []byte(strings.Split(stdout.String(), ":")[0])
				}

			}
			ids = append(ids, string(id))
		}
	}
	if len(ids) > 0 {
		return strings.Join(ids, ":"), nil
	}
	return "", nil
}

func GetEndpoint(ctx context.Context, ps process.Process) ([]string, error) {
	var (
		stdout    *bytes.Buffer
		err       error
		endpoints []string
	)

	stdout, err = ps.Run(
		ps.NewExecCommand(ctx, "nsenter", "-t", strconv.FormatInt(ps.Pid(), 10), "-n", "netstat", "-anp"),
	)
	if err != nil {
		return []string{}, err
	}
	stdout = command.Grep(stdout, `tcp|udp`,
		strconv.FormatInt(ps.Pid(), 10),
		`LISTEN`)
	endpointRaw := stdout.Bytes()
	for {
		if len(endpointRaw) == 0 {
			break
		}
		var val []byte
		protocols, _ := command.ReadField(endpointRaw, 1)

		val, endpointRaw = command.ReadField(endpointRaw, 4)
		if len(val) > 0 {
			endpoints = append(endpoints, string(protocols)+"/"+string(val))
		}
		endpointRaw = command.NextLine(endpointRaw)
	}

	return endpoints, nil

}

type ProcessBoard struct {
	total  int
	rwlock sync.RWMutex
}

func NewProcessBoard(total int) *ProcessBoard {
	return &ProcessBoard{
		total:  total,
		rwlock: sync.RWMutex{},
	}
}
func (p *ProcessBoard) Get() int {
	p.rwlock.RLock()
	defer p.rwlock.RUnlock()
	return p.total
}
func (p *ProcessBoard) Sub() {
	p.rwlock.Lock()
	defer p.rwlock.Unlock()
	p.total--
	return
}
