package core

import (
	"bytes"
	"context"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/Yeatesss/container-software/pkg/command"
	"github.com/Yeatesss/container-software/pkg/log"

	"github.com/Yeatesss/container-software/pkg/proc/process"
)

var _ SoftwareFinder = &CouchdbFindler{}

const Couchdb SwName = "couchdb"

type CouchdbFindler struct{}

func init() {
	if _, ok := Finders[DATABASE]; !ok {
		Finders[DATABASE] = make(map[SwName]SoftwareFinder)
	}
	Finders[DATABASE][Couchdb] = NewCouchdbFindler()
}
func NewCouchdbFindler() *CouchdbFindler {
	return &CouchdbFindler{}
}

type CouchdbResponse struct {
	Couchdb  string   `json:"couchdb"`
	Version  string   `json:"version"`
	GitSha   string   `json:"git_sha"`
	Uuid     string   `json:"uuid"`
	Features []string `json:"features"`
	Vendor   struct {
		Name string `json:"name"`
	} `json:"vendor"`
}

func (m CouchdbFindler) Verify(ctx context.Context, c *Container, thisis func(*Process, SoftwareFinder)) (bool, error) {
	var hit bool
	log.Logger.Debugf("Start verify influxDB:%s", c.Id)
	defer log.Logger.Debugf("Finish verify influxDB:%s", c.Id)

	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
		if ps._finder != nil {
			return
		}
		hit, err = m.SingleVerify(ctx, c.EnvPath, ps, thisis)
		return
	})
	return hit, err
}
func (m CouchdbFindler) SingleVerify(ctx context.Context, _ string, ps *Process, thisis func(*Process, SoftwareFinder)) (hit bool, err error) {
	var (
		eps             []string
		couchdbResponse CouchdbResponse
	)

	eps, err = GetEndpointWithChild(ctx, ps)
	if err != nil {
		return
	}
	for _, ep := range eps {
		var stdout *bytes.Buffer
		port := ep[strings.Index(ep, ":")+1:]
		cmds := append([]string{"-t", strconv.FormatInt(ps.Pid(), 10), "--pid", "--uts", "--ipc", "--net"}, []string{"curl", "http://localhost:" + port}...)
		cmd := ps.NewExecCommandWithEnv(ctx, "nsenter", cmds, "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
		stdout, err = ps.Run(cmd)
		if err != nil {
			return
		}
		err = jsoniter.UnmarshalFromString(stdout.String(), &couchdbResponse)
		if err != nil {
			return
		}
		if couchdbResponse.Version != "" && couchdbResponse.Couchdb != "" {
			ps.Version = couchdbResponse.Version
			hit = true
			thisis(ps, &m)
			return
		}
	}
	return
}
func (m CouchdbFindler) GetSoftware(ctx context.Context, c *Container) ([]*Software, error) {
	var softwares []*Software
	err := c.Processes.Range(func(_ int, ps *Process) (err error) {
		if ps._finder != nil {
			return nil
		}
		var software = &Software{
			Name:         "couchdb",
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
		software.BindEndpoint, err = GetEndpointWithChild(ctx, ps)
		if err != nil {
			return err
		}
		software.User, err = GetRunUser(ctx, ps, c.EnvPath)
		if err != nil {
			return err
		}
		software.Version = ps.Version

		software.ConfigPath, err = getCouchdbConfig(ctx, ps)
		if err != nil {
			return err
		}
		softwares = append(softwares, software)
		return nil
	})
	return softwares, err
}

func getCouchdbConfig(ctx context.Context, ps process.Process) (string, error) {
	var (
		stdout       *bytes.Buffer
		alternatives []string
		configRaw    []byte
		err          error
	)
	stdout, err = ps.Run(
		ps.EnterProcessNsRun(ctx, ps.Pid(), []string{"find", "/", "-path", "/proc", "-prune", "-o", "-path", "/lib", "-prune", "-o", "-path", "/lib64", "-prune", "-o", "-name", "etc", "-print"}),
	)
	if err != nil {
		return "", err
	}
	configRaw = stdout.Bytes()
	for len(configRaw) > 0 {
		raw := command.ReadLine(configRaw)
		if bytes.Contains(bytes.ToLower(raw), []byte("couchdb")) {
			alternatives = append(alternatives, string(raw))
		}
		configRaw = command.NextLine(configRaw)
	}
	return alternatives[0], nil
}
