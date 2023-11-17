package container_software

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Yeatesss/container-software/pkg/command"

	jsoniter "github.com/json-iterator/go"

	"github.com/Yeatesss/container-software/pkg/proc/process"

	"github.com/Yeatesss/container-software/core"
)

func TestUser(t *testing.T) {
	a, b := core.GetEndpoint(context.Background(), process.NewProcess(540824, nil))
	fmt.Println(a, b)
}
func TestFind(t *testing.T) {
	//stdout, err := process.NewProcess(12083, nil).Run(
	//	process.NewProcess(12083, nil).NewExecCommand(context.Background(), "bash", "-c", "nsenter -t 12083 -n && netstat -anp"),
	//)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(22305, stdout.String())
	mockMysqlCtr := &core.Container{
		Id:      "646d0c946f13f77eb8248b0c460bee6656f72397dd2db18ed634f13cac13de3d",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(4629, nil),
			},
		},
	}
	mockMysqlCtr.SetHypotheticalNspid()
	fmt.Println(mockMysqlCtr.Processes[0].NsPid())
	ctx := context.Background()

	find := NewFinder()

	f, e := find.Find(ctx, mockMysqlCtr)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)

}

func TestPostgrasqlFind(t *testing.T) {
	mockMysqlCtr := &core.Container{
		Id:      "646d0c946f13f77eb8248b0c460bee6656f72397dd2db18ed634f13cac13de3d",
		EnvPath: "PATH=/opt/bitnami/postgresql/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(58256, nil),
			},
			/*&core.Process{
				Process: process.NewProcess(2944380, nil),
			},*/
		},
	}
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockMysqlCtr, core.Postgresql)
	fmt.Println(jsoniter.MarshalToString(f))
	time.Sleep(10000 * time.Second)
	fmt.Println("sleep")
	fmt.Println(e)
}
func TestMongoFind(t *testing.T) {
	mockMongoCtr := &core.Container{
		Id:      "8b0fb3d01afc0d9f123a3487fffe14118214b99aad230443773eb6db0d99dc05",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(4091, nil),
			},
		},
	}
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockMongoCtr)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)
}

func TestSqlServerFind(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	mockMongoCtr := &core.Container{
		Id:      "46bfa47436ea8c238dcc6cc7335ab8b253d3ad34d388459579f3460b06dd1905",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Labels: map[string]string{
			"nerdctl/state-dir":                      "/var/lib/nerdctl/1935db59/containers/default/46bfa47436ea8c238dcc6cc7335ab8b253d3ad34d388459579f3460b06dd1905",
			"nerdctl/extraHosts":                     "null",
			"nerdctl/namespace":                      "default",
			"nerdctl/name":                           "sql-serve",
			"io.containerd.image.config.stop-signal": "SIGTERM",
			"containerd/namespaces":                  "default",
			"nerdctl/networks":                       "[\"bridge\"]",
			"nerdctl/hostname":                       "46bfa47436ea",
			"nerdctl/platform":                       "linux/amd64",
			"master_pid":                             "56944",
			"nerdctl/log-uri":                        "binary",
			"nerdctl/ports":                          "[{\"HostPort\"",
		},
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(56944, []int64{57167}),
			},
			&core.Process{
				Process: process.NewProcess(57167, []int64{}),
			},
		},
	}
	find := NewFinder()
	f, e := find.Find(ctx, mockMongoCtr)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)
}

func TestSqliteFind(t *testing.T) {
	mockSqliteCtr := &core.Container{
		Id:      "c90edaf54b53",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(119556, nil),
			},
		},
	}
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockSqliteCtr, core.Sqlite)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)
}

func TestRedisFind(t *testing.T) {
	mockRedisCtr := &core.Container{
		Id:      "8a3a59411794",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(6442, []int64{}),
			},
		},
	}
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockRedisCtr)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)
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
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockLighttpdCtr)
	ff, ee := find.Find(ctx, mockLighttpdCtr)
	fmt.Println(f, e, ff, ee)
}
func TestNginxFind(t *testing.T) {
	//{"Process":{"Pid":28157,"ChildPids":[]},"version":""},{"Process":{"Pid":28158,"ChildPids":[]},"version":""},{"Process":{"Pid":28159,"ChildPids":[]},"version":""},{"Process":{"Pid":28160,"ChildPids":[]},"version":""}]}`
	mockNginxCtr := &core.Container{
		Id:      "509d203eca13972c071c2530b9b6162956186faa3a05656dd6dbc202352aab6a",
		EnvPath: "",
		Labels:  map[string]string{},
		Processes: []*core.Process{
			{
				Process: process.NewProcess(28112, []int64{28158, 28159, 28157, 28160}),
				Version: "",
			},
			{
				Process: process.NewProcess(28157, []int64{}),
				Version: "",
			},
			{
				Process: process.NewProcess(28158, []int64{}),
				Version: "",
			},
			{
				Process: process.NewProcess(28159, []int64{}),
				Version: "",
			},
			{
				Process: process.NewProcess(28160, []int64{}),
				Version: "",
			},
		},
	}
	//fmt.Println(jsoniter.UnmarshalFromString(a, mockNginxCtr))
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockNginxCtr)
	//ff, ee := find.Find(ctx, mockNginxCtr)
	fmt.Println(f, e)
}

func TestTomcatFind(t *testing.T) {

	mockLighttpdCtr := &core.Container{
		Id:      "a5ca5fd4cfd3e74b2eb6dc85325a8192d2955ccb2890bdbe14ef5021e905cbf5",
		EnvPath: "",
		Labels: map[string]string{
			"nerdctl/namespace":                      "default",
			"nerdctl/hostname":                       "a5ca5fd4cfd3",
			"nerdctl/platform":                       "linux/amd64",
			"nerdctl/log-uri":                        "binary",
			"nerdctl/anonymous-volumes":              "[\"959b7a8704a088a318d00877adbf7232de7cb7775c6cf791e976005c6f17f344\"]",
			"nerdctl/extraHosts":                     "null",
			"nerdctl/state-dir":                      "/var/lib/nerdctl/1935db59/containers/default/a5ca5fd4cfd3e74b2eb6dc85325a8192d2955ccb2890bdbe14ef5021e905cbf5",
			"containerd/namespaces":                  "default",
			"nerdctl/networks":                       "[\"bridge\"]",
			"master_pid":                             "92521",
			"io.containerd.image.config.stop-signal": "SIGTERM",
			"nerdctl/mounts":                         "[{\"Type\"",
			"nerdctl/name":                           "mysql-container",
		},
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(92521, []int64{}),
			},
		},
	}
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockLighttpdCtr)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(f, e)

}

func TestJbossFind(t *testing.T) {
	mockJbossCtr := &core.Container{
		Id:      "1d209baf766fed83dfc3ec2d740bd940d18f788258ffbff6cd20b00b9c149a45",
		EnvPath: "PATH=/opt/java/openjdk/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(101572, []int64{101805}),
			},
			&core.Process{
				Process: process.NewProcess(101805, nil),
			},
		},
	}
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockJbossCtr)
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
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockApacheCtr)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)
}
func TestFindPOD(t *testing.T) {

	mockPODCtr := core.Container{
		Id:      "00d46a62ae2bf5ddfa4cc0078dcea75cda088fabb943fa8e398ab5028e5d3914",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
	}

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	ctx, _ := context.WithTimeout(context.Background(), 9999*time.Second)
	//for i := 0; i <= 100; i++ {
	//	eg.Go(func() error {
	c := mockPODCtr
	c.Processes = core.Processes{
		&core.Process{
			Process: process.NewProcess(93881, nil),
		},
		//&core.Process{
		//	Process: process.NewProcess(6899, nil),
		//},
		//&core.Process{
		//	Process: process.NewProcess(55829, nil),
		//},
	}
	find := NewFinder()
	f, e := find.Find(ctx, &c)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)

	// 创建一个管道来捕获命令输出
	//cmd := exec.CommandContext(ctx, "nsenter", "-t",
	//	"5279",
	//	"--pid",
	//	"--uts",
	//	"--ipc",
	//	"--net",
	//	"--mount",
	//	"/pause",
	//	"-help")
	//stdout, err := cmd.StdoutPipe()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//stderr, err := cmd.StderrPipe()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// 启动命令
	//if err := cmd.Start(); err != nil {
	//	fmt.Println(err)
	//}
	//go func() {
	//	var buf bytes.Buffer
	//	var buf1 bytes.Buffer
	//	_, err = io.Copy(&buf1, stderr)
	//	if err != nil {
	//		fmt.Println("1111", err, buf1.String())
	//	}
	//	_, err = io.Copy(&buf, stdout)
	//	if err != nil {
	//		fmt.Println("2222", err, buf.String())
	//	}
	//	// 处理命令输出
	//	output := buf.String()
	//	fmt.Println(output)
	//}()
	//time.Sleep(20 * time.Second)
	// 将命令输出复制到一个缓冲区

	//go func() {
	//	time.Sleep(6 * time.Second)
	//	fmt.Println("sle")
	//	fmt.Println(cmd.Process.Kill())
	//}()
	// 等待命令执行完成

	//if err := cmd.Wait(); err != nil {
	//	fmt.Println(2222)
	//
	//	log.Fatal(err)
	//}
	//fmt.Println(111)
	//d, e := exec.Command("nsenter", "-t", "55828", "--pid", "--uts", "--ipc", "--net", "--mount", "getent", "passwd").CombinedOutput()
	//fmt.Println(string(d), e)
	//aaa := bytes.NewBuffer([]byte(`
	//root:x:0:0:root:/root:/bin/bash
	//daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin
	//bin:x:2:2:bin:/bin:/usr/sbin/nologin
	//sys:x:3:3:sys:/dev:/usr/sbin/nologin
	//sync:x:4:65534:sync:/bin:/bin/sync
	//games:x:5:60:games:/usr/games:/usr/sbin/nologin
	//man:x:6:12:man:/var/cache/man:/usr/sbin/nologin
	//lp:x:7:7:lp:/var/spool/lpd:/usr/sbin/nologin
	//mail:x:8:8:mail:/var/mail:/usr/sbin/nologin
	//news:x:9:9:news:/var/spool/news:/usr/sbin/nologin
	//uucp:x:10:10:uucp:/var/spool/uucp:/usr/sbin/nologin
	//proxy:x:13:13:proxy:/bin:/usr/sbin/nologin
	//www-data:x:33:33:www-data:/var/www:/usr/sbin/nologin
	//backup:x:34:34:backup:/var/backups:/usr/sbin/nologin
	//list:x:38:38:Mailing List Manager:/var/list:/usr/sbin/nologin
	//irc:x:39:39:ircd:/run/ircd:/usr/sbin/nologin
	//_apt:x:42:65534::/nonexistent:/usr/sbin/nologin
	//nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin
	//nginx:x:101:101:nginx user:/nonexistent:/bin/false
	//`))
	//var bufOut = &bytes.Buffer{}
	//
	//cmd := exec.Command("grep", "101")
	//sp, e1 := cmd.StdinPipe()
	//fmt.Println(e1)
	//io.Copy(sp, aaa)
	////cmd.Stdin = aaa
	//o, e := cmd.StdoutPipe()
	//
	//err := cmd.Start() // 改成 Start 而不是 Run
	//if err != nil {
	//	return
	//}
	//fmt.Println(9090)
	//fmt.Println(e)
	//fmt.Println(io.Copy(bufOut, o))
	//fmt.Println(bufOut.String())
	//if err = cmd.Wait(); err != nil {
	//	fmt.Println(1232312)
	//}
	//return
	//runner := command.NewCmdRuner()
	//a, b := runner.Run(
	//	runner.NewExecCommand(context.Background(), "nsenter", "-t", "55724", "--pid", "--uts", "--ipc", "--net", "--mount", "/usr/sbin/nginx", "-V"),
	//	//runner.NewExecCommand(context.Background(), "xargs", "echo"),
	//)
	//fmt.Println(a.String(), b)
	////var buf = &bytes.Buffer{}
	////cmd.Stdout = buf
	//go func() {
	//	time.Sleep(5 * time.Second)
	//	fmt.Println(cmd.Process.Pid)
	//	fmt.Println(cmd.Process.Kill())
	//}()
	//stdoutPipe, e := cmd.CombinedOutput()
	//
	//fmt.Println(string(stdoutPipe), e)
	////go func() {
	////	time.Sleep(6 * time.Second)
	////	fmt.Println("stop")
	////	out := &bytes.Buffer{}
	////	fmt.Println(io.Copy(out, stdoutPipe))
	////	fmt.Println(out.String())
	////	stdoutPipe.Close()
	////}()
	//////
	////io.Copy(os.Stdout)
	////fmt.Println(cmd.Run())
	//time.Sleep(8 * time.Second)
	//fmt.Println(buf.String())

}
func TestCsds(t *testing.T) {

}

func Ttt() {
	mockMongoCtr := &core.Container{
		Id:      "329c117ebfe813efd66bd60cd71d7b4f3bb474a5375228f92ae902e312c292da",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Labels: map[string]string{
			"annotation.io.kubernetes.container.terminationMessagePolicy": "File",
			"org.label-schema.name":                                     "CentOS Base Image",
			"master_pid":                                                "9866",
			"org.opencontainers.image.title":                            "CentOS Base Image",
			"annotation.io.kubernetes.container.ports":                  "[{\"containerPort\"",
			"org.label-schema.build-date":                               "20201113",
			"org.label-schema.license":                                  "GPLv2",
			"io.kubernetes.container.name":                              "ngep-test",
			"annotation.io.kubernetes.pod.terminationGracePeriod":       "30",
			"org.opencontainers.image.created":                          "2020-11-13 00",
			"annotation.io.kubernetes.container.hash":                   "b5f89fad",
			"annotation.io.kubernetes.container.terminationMessagePath": "/dev/termination-log",
			"io.kubernetes.container.logpath":                           "/var/log/pods/default_ngep-manager-5dd7cf4df5-s8572_d4bd40ca-72f1-4e9a-b862-9ac57d10b024/ngep-test/4.log",
			"org.label-schema.vendor":                                   "CentOS",
			"io.kubernetes.pod.namespace":                               "default",
			"org.opencontainers.image.licenses":                         "GPL-2.0-only",
			"org.opencontainers.image.vendor":                           "CentOS",
			"org.label-schema.schema-version":                           "1.0",
			"io.kubernetes.pod.uid":                                     "d4bd40ca-72f1-4e9a-b862-9ac57d10b024",
			"io.kubernetes.pod.name":                                    "ngep-manager-5dd7cf4df5-s8572",
			"io.kubernetes.sandbox.id":                                  "cf0789ca75cabeacdd084c826de7ca86b8ba9bfc3f5a3122795139d021595d60",
			"annotation.io.kubernetes.container.restartCount":           "4",
			"io.kubernetes.docker.type":                                 "container",
		},
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(9866, []int64{11821, 10010, 10059, 10060, 10109, 10110}),
			},

			&core.Process{
				Process: process.NewProcess(10010, nil),
			},
			&core.Process{
				Process: process.NewProcess(10059, nil),
			},
			&core.Process{
				Process: process.NewProcess(10060, nil),
			},
			&core.Process{
				Process: process.NewProcess(10109, nil),
			},
			&core.Process{
				Process: process.NewProcess(10110, nil),
			},
			&core.Process{
				Process: process.NewProcess(11821, []int64{11889, 11847, 11879, 23323, 11845, 26022, 11842, 11838, 11843, 11840, 11841, 11839}),
			},
			&core.Process{
				Process: process.NewProcess(11838, nil),
			},
			&core.Process{
				Process: process.NewProcess(11839, nil),
			},
			&core.Process{
				Process: process.NewProcess(11840, nil),
			},
			&core.Process{
				Process: process.NewProcess(11841, []int64{15819}),
			},
			&core.Process{
				Process: process.NewProcess(11842, nil),
			},
			&core.Process{
				Process: process.NewProcess(11843, nil),
			},
			&core.Process{
				Process: process.NewProcess(11845, []int64{12611}),
			},
			&core.Process{
				Process: process.NewProcess(11847, nil),
			},
			&core.Process{
				Process: process.NewProcess(11879, nil),
			},
			&core.Process{
				Process: process.NewProcess(11889, []int64{69910, 69909, 69911, 69912, 69915, 69914, 69913, 69916}),
			},
			&core.Process{
				Process: process.NewProcess(12611, []int64{14238}),
			},
			&core.Process{
				Process: process.NewProcess(14238, nil),
			},
			&core.Process{
				Process: process.NewProcess(15819, nil),
			},
			&core.Process{
				Process: process.NewProcess(23323, nil),
			},
			&core.Process{
				Process: process.NewProcess(26022, nil),
			},
			&core.Process{
				Process: process.NewProcess(69909, nil),
			},
			&core.Process{
				Process: process.NewProcess(69910, nil),
			},
			&core.Process{
				Process: process.NewProcess(69911, nil),
			},
			&core.Process{
				Process: process.NewProcess(69912, nil),
			},
			&core.Process{
				Process: process.NewProcess(69913, nil),
			},
			&core.Process{
				Process: process.NewProcess(69914, nil),
			},
			&core.Process{
				Process: process.NewProcess(69915, nil),
			},
			&core.Process{
				Process: process.NewProcess(69916, nil),
			},
		},
	}
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	find := NewFinder()
	f, e := find.Find(ctx, mockMongoCtr)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)
}

func TestNgepFind(t *testing.T) {
	//c, d := exec.Command("nsenter", "-t", "107751", "--pid", "--uts", "--ipc", "--net", "--mount", "cat", "/proc/49/cmdline").CombinedOutput()
	//fmt.Println(c, d)
	cmdline, err := command.NewCmdRuner().Run(
		command.NewCmdRuner().EnterProcessNsRun(context.Background(), 107751, []string{"cat", "/proc/" + "49" + "/cmdline"}, "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"),
	)
	fmt.Println(cmdline, err)
	cmdline, err = command.NewCmdRuner().Run(
		command.NewCmdRuner().EnterProcessNsRun(context.Background(), 107751, []string{"cat", "/proc/" + "49" + "/cmdline"}, "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"),
	)
	fmt.Println(cmdline, err)
	//core.GetEndpoint(context.Background(), process.NewProcess(11879, nil))
	//var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	Ttt()
	//}()
	//
	//wg.Wait()

}

func TestInfluxDB(t *testing.T) {
	mockSqliteCtr := &core.Container{
		Id:      "5201ff905b57",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(387088, nil),
			},
		},
	}
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockSqliteCtr, core.InfluxDB)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)
}
func BenchmarkCouchdb(b *testing.B) {
	mockSqliteCtr := &core.Container{
		Id:      "5201ff905b57",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(388032, []int64{388061}),
			},
			&core.Process{
				Process: process.NewProcess(388061, nil),
			},
		},
	}
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockSqliteCtr, core.Couchdb)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)
}

func TestCouchdb(t *testing.T) {
	mockSqliteCtr := &core.Container{
		Id:      "5201ff905b57",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(388032, []int64{388061}),
			},
			&core.Process{
				Process: process.NewProcess(388061, nil),
			},
		},
	}
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockSqliteCtr)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)
}
