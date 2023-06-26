package container_software

import (
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
		ctrCache = freecache.NewCache(1024 * 1024 * 20)
	})
	return &Finder{Cache{ctrCache: ctrCache}}
}

func (l *Finder) Find(c *core.Container, onlyTypes ...core.SwType) (softwares []*core.Software, err error) {
	var (
		hit                  bool
		totalNumberProcesses = c.Processes.Len()
	)
	if val, _ := l.ctrCache.Get([]byte(c.Id)); len(val) > 0 {
		err = jsoniter.Unmarshal(val, &softwares)
		return
	}
	var validContainer = func(finders map[string]core.SoftwareFinder) {
		if totalNumberProcesses == 0 {
			return
		}
		for _, finder := range finders {
			if totalNumberProcesses == 0 {
				return
			}
			if finder.Verify(c, func(p *core.Process, finder core.SoftwareFinder) {
				p.SetFinder(finder)
			}) {
				totalNumberProcesses--
				hit = true
			}
		}
		return
	}
	if len(onlyTypes) > 0 {
		for _, onlyType := range onlyTypes {
			validContainer(core.Finders[onlyType])
		}
	} else {
		for _, finders := range core.Finders {
			validContainer(finders)
		}
	}
	if hit {
		var softwaresByte []byte
		softwares, err = core.GetSoftware(c)
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
