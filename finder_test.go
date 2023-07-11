package container_software

import (
	"fmt"
	"testing"

	jsoniter "github.com/json-iterator/go"

	"github.com/Yeatesss/container-software/pkg/proc/process"

	"github.com/Yeatesss/container-software/core"
)

func TestFind(t *testing.T) {
	mockMysqlCtr := &core.Container{
		Id:      "0e16a4e3751d8df3b70f26957c9a6fb72405e89c6aa7820186acb96a79298af3",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(313838, nil),
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
				Process: process.NewProcess(2548908, nil),
			},
		},
	}
	find := NewFinder()
	f, e := find.Find(mockLighttpdCtr)
	ff, ee := find.Find(mockLighttpdCtr)
	fmt.Println(f, e, ff, ee)
}
func TestNginxFind(t *testing.T) {
	mockNginxCtr := &core.Container{
		Id:      "d58f8a7897d5aac5c6c0817204e765f9698e4b4266b1ac536b507e021a4da54d",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			//&core.Process{
			//	Process: process.NewProcess(1929787,nil),
			//},
			&core.Process{
				Process: process.NewProcess(1929934, nil),
			},
		},
	}
	find := NewFinder()
	f, e := find.Find(mockNginxCtr)
	ff, ee := find.Find(mockNginxCtr)
	fmt.Println(f, e, ff, ee)
}
func TestTomcatFind(t *testing.T) {
	mockLighttpdCtr := &core.Container{
		Id:      "ab7920723bbfe58198076216a43fffe0e5a2e88c04c46f7bc6611bb0c8a25a04",
		EnvPath: "PATH=/usr/local/tomcat/bin:/opt/java/openjdk/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(1142497, nil),
			},
		},
	}
	find := NewFinder()
	f, e := find.Find(mockLighttpdCtr)
	ff, ee := find.Find(mockLighttpdCtr)
	fmt.Println(f, e, ff, ee)
}

func TestJbossFind(t *testing.T) {
	mockJbossCtr := &core.Container{
		Id:      "5c95bddc2a3c8a4e94b58ecda66564eb32d7e314985199fb6413a6b32feeca21",
		EnvPath: "PATH=/opt/java/openjdk/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(2074003, []int64{2074133}),
			},
			&core.Process{
				Process: process.NewProcess(2074133, nil),
			},
		},
	}
	find := NewFinder()
	f, e := find.Find(mockJbossCtr)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)
}

func TestApacheFind(t *testing.T) {
	mockApacheCtr := &core.Container{
		Id:      "ca77d884536f6f7faec7662efb7ec652ec5627edeeb910d7ac50ae62e95f04e4",
		EnvPath: "PATH=/usr/local/apache2/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(2549367, []int64{2549399}),
			},
			&core.Process{
				Process: process.NewProcess(2549399, []int64{2549403}),
			},
		},
	}
	find := NewFinder()
	f, e := find.Find(mockApacheCtr)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)
}
