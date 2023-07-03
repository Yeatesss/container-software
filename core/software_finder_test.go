package core

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

func TestFinders(t *testing.T) {
	mockMysqlCtr := &Container{
		Id:      "f48d3f77afc75c526fddee8de075e5cf7c25a6cf038febbba6a09b6d4da6b421",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: Processes{
			&Process{
				Process: process.NewProcess(313838, nil),
			},
			//&Process{
			//	Pid: 1572497,
			//	Exe: "lrwxrwxrwx 1 lxd docker 0 Jun 15 15:00 /proc/1572497/exe -> /usr/sbin/mysqld",
			//	Cmd: "mysqld --daemonize --skip-networking --default-time-zone=SYSTEM --socket=/var/run/mysqld/mysqld.sock ",
			//},
		},
	}
	fmt.Printf("%p\n", Finders[DATABASE]["mysql"])
	software, err := Finders[DATABASE]["mysql"].GetSoftware(mockMysqlCtr)
	if err != nil {
		return
	}
	fmt.Println(software)
}

func TestRegexp(t *testing.T) {
	str := "--defaults-file = /path/to/file"

	re := regexp.MustCompile(`--defaults-file[\x20=]+(\S+)`)
	match := re.FindStringSubmatch(str)

	if len(match) > 1 {
		content := match[1]
		fmt.Println("括号中的内容:", content)
	} else {
		fmt.Println("未找到匹配的内容")
	}
}
