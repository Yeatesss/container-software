package core

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/Yeatesss/container-software/pkg/command"
	"github.com/Yeatesss/container-software/pkg/log"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

var _ SoftwareFinder = &InfluxDBFindler{}

const InfluxDB SwName = "influxDB"

type InfluxDBFindler struct{}

func init() {
	if _, ok := Finders[DATABASE]; !ok {
		Finders[DATABASE] = make(map[SwName]SoftwareFinder)
	}
	Finders[DATABASE][InfluxDB] = NewInfluxDBFindler()
}
func NewInfluxDBFindler() *InfluxDBFindler {
	return &InfluxDBFindler{}
}

func (m InfluxDBFindler) Verify(ctx context.Context, c *Container, thisis func(*Process, SoftwareFinder)) (bool, error) {
	var hit bool
	log.Logger.Debugf("Start verify influxDB:%s", c.Id)
	defer log.Logger.Debugf("Finish verify influxDB:%s", c.Id)

	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
		if ps._finder != nil {
			return nil
		}
		//fmt.Println(c.Id, ps.Pid(), ps.ChildPids())
		var exe string
		exe, err = process.GetProcessExe(ctx, ps.Process)
		if err != nil {
			return
		}
		cmd := ps.EnterProcessNsRun(ctx, ps.Pid(), []string{exe, "version"}, "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
		stdout, err := ps.Run(cmd)
		if err != nil {
			//TODO:Don't know why it returns a correct result and also returns an err
			if _, ok := err.(*exec.ExitError); ok {
				err = nil
			} else {
				return err
			}
		}
		if strings.Contains(strings.ToLower(stdout.String()), "influxd") {
			versionByte, _ := command.ReadField(stdout.Bytes(), 2)
			ps.Version = string(versionByte)
			hit = true
			thisis(ps, &m)
			return nil
		}
		return
	})
	return hit, err
}

func (m InfluxDBFindler) GetSoftware(ctx context.Context, c *Container) ([]*Software, error) {
	var softwares []*Software
	err := c.Processes.Range(func(_ int, ps *Process) (err error) {

		var software = &Software{
			Name:         "influxDB",
			Type:         DATABASE,
			Version:      "",
			BindEndpoint: nil,
			User:         "",
			BinaryPath:   "",
			ConfigPath:   "",
		}
		var exe string
		exe, err = process.GetProcessExe(ctx, ps.Process)
		if err != nil {
			return
		}
		software.BinaryPath = exe
		software.BindEndpoint, err = GetEndpoint(ctx, ps)
		if err != nil {
			return err
		}
		software.User, err = GetRunUser(ctx, ps, c.EnvPath)
		if err != nil {
			return err
		}
		if ps.Version != "" {
			software.Version = ps.Version
		} else {
			software.Version, err = getInfluxDBVersion(ctx, ps, exe)
			if err != nil {
				return err
			}
		}

		software.ConfigPath, err = getInfluxDBConfig(ctx, ps)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, err
}

func getInfluxDBConfig(ctx context.Context, ps process.Process) (string, error) {
	cmdline, err := ps.Cmdline()
	if err != nil {
		return "-", err
	}
	cmdlineStr := cmdline.String()
	if strings.Contains(cmdlineStr, "config") {
		cmdlineStr = strings.ReplaceAll(cmdlineStr, "\u0000", " ")
		cmdlineBytes := []byte(cmdlineStr)
		for len(cmdlineBytes) > 0 {
			var hit []byte
			hit, cmdlineBytes = command.NextField(cmdlineBytes)
			if bytes.Contains(hit, []byte("-config")) {
				hit, cmdlineBytes = command.NextField(cmdlineBytes)
				return string(hit), nil
			}
		}
	}

	//command.ReadField(cmdline, 2)
	return "-", nil
}
func getInfluxDBVersion(ctx context.Context, ps process.Process, exe string) (string, error) {
	var (
		stdout *bytes.Buffer
		err    error
	)
	ps.CacheClear(ps.EnterProcessNsRun(ctx, ps.Pid(), []string{exe, "-v"}))

	stdout, err = ps.Run(
		ps.EnterProcessNsRun(ctx, ps.Pid(), []string{exe, "-v"}),
	)
	if err != nil {
		return "", err
	}
	val, _ := command.ReadField(stdout.Bytes(), 8)
	return string(val), nil

}

//
//func GetSqlserverEndpoint(ctx context.Context, ps process.Process) ([]string, error) {
//	var (
//		stdout    *bytes.Buffer
//		err       error
//		endpoints []string
//		pids      []string
//	)
//
//	stdout, err = ps.Run(
//		ps.NewExecCommand(ctx, "nsenter", "-t", strconv.FormatInt(ps.Pid(), 10), "-n", "netstat", "-anp"),
//	)
//	if err != nil {
//		return []string{}, err
//	}
//	pids = append(pids, strconv.FormatInt(ps.Pid(), 10))
//	for _, pid := range ps.ChildPids() {
//		pids = append(pids, strconv.FormatInt(pid, 10))
//
//	}
//	stdout = command.Grep(stdout, `tcp>>LISTEN|udp`,
//		strings.Join(pids, "|"))
//	endpointRaw := stdout.Bytes()
//	for {
//		if len(endpointRaw) == 0 {
//			break
//		}
//		var val []byte
//		protocols, _ := command.ReadField(endpointRaw, 1)
//
//		val, endpointRaw = command.ReadField(endpointRaw, 4)
//		if len(val) > 0 {
//			endpoints = append(endpoints, string(protocols)+"/"+string(val))
//		}
//		endpointRaw = command.NextLine(endpointRaw)
//	}
//
//	return endpoints, nil
//
//}
