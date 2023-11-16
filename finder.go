package container_software

import (
	"context"
	"regexp"
	"strings"
	"sync"

	"github.com/shirou/gopsutil/v3/process"

	"github.com/pkg/errors"

	jsoniter "github.com/json-iterator/go"

	"github.com/Yeatesss/container-software/core"

	"github.com/coocood/freecache"
)

var (
	ctrCache = freecache.NewCache(1024 * 1024 * 10)
)

type Finder struct {
	Cache
}
type Cache struct {
	ctrCache *freecache.Cache
}

func NewFinder() *Finder {
	return &Finder{Cache{ctrCache: ctrCache}}
}

func (l *Finder) Find(ctx context.Context, c *core.Container, onlys ...interface{}) (softwares []*core.Software, err error) {
	var (
		hit          bool
		processBoard = core.NewProcessBoard(c.Processes.Len())
		wg           sync.WaitGroup
	)
	//set environment variable
	if c.EnvPath == "" {
		ps, _ := process.NewProcess(int32(c.Processes[0].Process.Pid()))
		psenv, _ := ps.Environ()
		for _, env := range psenv {
			if strings.HasPrefix(env, "PATH") {
				c.EnvPath = env
				break
			}
		}
	}
	if val, _ := l.ctrCache.Get([]byte(c.Id)); len(val) > 0 {
		err = jsoniter.Unmarshal(val, &softwares)
		return
	}

	var validContainer = func(finder core.SoftwareFinder) {
		var err error
		if processBoard.Get() == 0 {
			return
		}
		var is bool
		if processBoard.Get() == 0 {
			return
		}
		if is, err = finder.Verify(ctx, c, func(p *core.Process, finder core.SoftwareFinder) {
			p.SetFinder(finder)
		}); err == nil && is {
			processBoard.Sub()
			hit = true
		}

		return
	}
	var validContainers = func(finders map[core.SwName]core.SoftwareFinder) {
		if processBoard.Get() == 0 {
			return
		}
		for _, finder := range finders {
			wg.Add(1)
			go func(finder core.SoftwareFinder) {
				defer wg.Done()
				validContainer(finder)
			}(finder)
		}
		wg.Wait()
		return
	}
	var priorityMarking = func(envPath string, ps *core.Process) error {
		var is bool
		cmdline, err := ps.Cmdline()
		if err != nil {
			return err
		}
		if strings.Contains(strings.ToLower(cmdline.String()), "mysql") {
			if is, err = core.NewMysqlFindler().SingleVerify(ctx, ps, func(p *core.Process, finder core.SoftwareFinder) {
				p.SetFinder(finder)
			}); err == nil && is {
				hit = true
			}
		}
		if strings.Contains(strings.ToLower(cmdline.String()), "sqlservr") {
			if is, err = core.NewSqlServerFindler().SingleVerify(ctx, ps, func(p *core.Process, finder core.SoftwareFinder) {
				p.SetFinder(finder)
			}); err == nil && is {
				hit = true
			}
		}
		if strings.Contains(strings.ToLower(cmdline.String()), "jboss") || strings.Contains(strings.ToLower(cmdline.String()), "wildfly") {
			if is, err = core.NewJbossFindler().SingleVerify(ctx, envPath, ps, func(p *core.Process, finder core.SoftwareFinder) {
				p.SetFinder(finder)
			}); err == nil && is {
				hit = true
			}
		}
		return nil
	}

	if len(onlys) > 0 {
		for _, only := range onlys {
			switch v := only.(type) {
			case core.SwType:
				validContainers(core.Finders[v])
			case core.SwName:
				for _, finders := range core.Finders {
					if finder, ok := finders[v]; ok {
						validContainer(finder)
					}
				}
			}
		}
	} else {
		for _, process := range c.Processes {
			_ = priorityMarking(c.EnvPath, process)
		}

		for _, finders := range core.Finders {
			validContainers(finders)
		}
	}
	if hit {
		var softwaresByte []byte
		softwares, err = core.GetSoftware(ctx, c)
		if err != nil {
			return
		}
		err = Check(softwares)
		if err != nil {
			return
		}
		softwaresByte, err = jsoniter.Marshal(softwares)
		if err != nil {
			return
		}
		l.ctrCache.Set([]byte(c.Id), softwaresByte, 0)
	}
	return
}

func Check(sfs []*core.Software) error {
	pattern := regexp.MustCompile(`[vV]*\d+\.\S+`)

	for _, sf := range sfs {
		if !pattern.MatchString(sf.Version) {
			return errors.New(sf.Name + "-version " + sf.Version + " is not valid")
		}
		if strings.Contains(sf.ConfigPath, "find:") {
			return errors.New(sf.Name + "-config " + sf.ConfigPath + " is not valid")
		}
	}
	return nil
}
