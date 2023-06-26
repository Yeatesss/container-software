package container_software

import (
	"fmt"
	"testing"

	"github.com/Yeatesss/container-software/pkg/proc/process"

	"github.com/Yeatesss/container-software/core"
)

func TestFind(t *testing.T) {
	mockMysqlCtr := &core.Container{
		Id:      "0e16a4e3751d8df3b70f26957c9a6fb72405e89c6aa7820186acb96a79298af3",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(313838),
			},
		},
	}
	find := NewFinder()
	find.Find(mockMysqlCtr)
}

func TestLighttpdFind(t *testing.T) {
	mockLighttpdCtr := &core.Container{
		Id:      "1f189cefd12834a0e40574e951f767bbfe5de961172584f8fba5ef876679402d",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(2548908),
			},
		},
	}
	find := NewFinder()
	f, e := find.Find(mockLighttpdCtr)
	ff, ee := find.Find(mockLighttpdCtr)
	fmt.Println(f, e, ff, ee)
}

func TestTomcatFind(t *testing.T) {
	mockLighttpdCtr := &core.Container{
		Id:      "ab7920723bbfe58198076216a43fffe0e5a2e88c04c46f7bc6611bb0c8a25a04",
		EnvPath: "PATH=/usr/local/tomcat/bin:/opt/java/openjdk/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(1142497),
			},
		},
	}
	find := NewFinder()
	f, e := find.Find(mockLighttpdCtr)
	ff, ee := find.Find(mockLighttpdCtr)
	fmt.Println(f, e, ff, ee)
}

//func TestByte(t *testing.T) {
//	var b = []byte{9}
//	fmt.Println(string(b))
//}
