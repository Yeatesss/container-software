package container_software

import (
	"context"
	"sync"

	jsoniter "github.com/json-iterator/go"

	"github.com/Yeatesss/container-software/core"

	"github.com/coocood/freecache"
)

var (
	ctrCache     *freecache.Cache
	ctrCacheOnce sync.Once
)

type Finder struct {
	Cache
}
type Cache struct {
	ctrCache *freecache.Cache
}

func NewFinder() *Finder {
	ctrCacheOnce.Do(func() {
		ctrCache = freecache.NewCache(1024 * 1024 * 10)
	})
	return &Finder{Cache{ctrCache: ctrCache}}
}

func (l *Finder) Find(ctx context.Context, c *core.Container, onlys ...interface{}) (softwares []*core.Software, err error) {
	var (
		hit          bool
		processBoard = core.NewProcessBoard(c.Processes.Len())
		wg           sync.WaitGroup
	)
	if val, _ := l.ctrCache.Get([]byte(c.Id)); len(val) > 0 {
		err = jsoniter.Unmarshal(val, &softwares)
		return
	}
	var validContainer = func(finder core.SoftwareFinder) {
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
		softwaresByte, err = jsoniter.Marshal(softwares)
		if err != nil {
			return
		}
		l.ctrCache.Set([]byte(c.Id), softwaresByte, 0)
	}
	return
}
