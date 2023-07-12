package container_software

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

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
	ctx := context.Background()
	find := NewFinder()
	find.Find(ctx, mockMysqlCtr)
}

func TestPostgrasqlFind(t *testing.T) {
	mockMysqlCtr := &core.Container{
		Id:      "9cf0090ccb274d9702d0abc229572a27b3a9b5646994df64030f7dd0ccd4a4cc",
		EnvPath: "PATH=/opt/bitnami/postgresql/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(2944176, []int64{2944380}),
			},
			&core.Process{
				Process: process.NewProcess(2944380, nil),
			},
		},
	}
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockMysqlCtr)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)
}
func TestMongoFind(t *testing.T) {
	mockMongoCtr := &core.Container{
		Id:      "8b0fb3d01afc0d9f123a3487fffe14118214b99aad230443773eb6db0d99dc05",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(3074901, []int64{3074990}),
			},
			&core.Process{
				Process: process.NewProcess(3074990, nil),
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
	mockMongoCtr := &core.Container{
		Id:      "7f20fb5e99fba9377480765e2249d0d859b5ef69d62f25c193d8f7a7ed08543a",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(3082182, []int64{3082239}),
			},
			&core.Process{
				Process: process.NewProcess(3082239, []int64{3082241}),
			},
			&core.Process{
				Process: process.NewProcess(3082241, nil),
			},
		},
	}
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockMongoCtr)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)
}

func TestSqliteFind(t *testing.T) {
	mockSqliteCtr := &core.Container{
		Id:      "44bb4e0f6243534cf97d54f63c3b527aaf339ca245209c592230580b8240e28d",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(3878895, nil),
			},
		},
	}
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockSqliteCtr)
	fmt.Println(jsoniter.MarshalToString(f))
	fmt.Println(e)
}

func TestRedisFind(t *testing.T) {
	mockRedisCtr := &core.Container{
		Id:      "2689c2fff772472bde16f59817ad297acd2df0b234a660a487d97ee1c13fb1a2",
		EnvPath: "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		Processes: core.Processes{
			&core.Process{
				Process: process.NewProcess(3946496, []int64{3946536}),
			},
			&core.Process{
				Process: process.NewProcess(3946536, nil),
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
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockNginxCtr)
	ff, ee := find.Find(ctx, mockNginxCtr)
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
	ctx := context.Background()
	find := NewFinder()
	f, e := find.Find(ctx, mockLighttpdCtr)
	ff, ee := find.Find(ctx, mockLighttpdCtr)
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
	c, b := exec.Command("go", "env"), new(bytes.Buffer)
	c.Stdout = b
	c.Run()

	s := bufio.NewScanner(b)
	for s.Scan() {
		if strings.Contains(s.Text(), "CACHE") {
			println(s.Text())
		}
	}
}

//stdout, err := core.Process{
//Process: process.NewProcess(5279, nil),
//}.Run(
//core.Process{
//Process: process.NewProcess(5279, nil),
//}.NewExecCommand(ctx, "nsenter", "-t",
//"5279",
//"--pid",
//"--uts",
//"--ipc",
//"--net",
//"--mount",
//"/pause",
//"-v"))
//fmt.Println(stdout)
//fmt.Println(err)
