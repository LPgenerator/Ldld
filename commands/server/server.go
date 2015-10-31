package server

import (
	"os"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/LPgenerator/Ldld/helpers"
	"github.com/LPgenerator/Ldld/helpers/backends"
	"github.com/LPgenerator/Ldld/helpers/providers"
)

type LdlSrv struct{
	path      string
	dist      string
	cli_path  string
	fs        string
	backend   backends.Fs
	provider  providers.Provider
}

func New(path string, dist string, cli_path string, fs string) (*LdlSrv) {
	strm := &LdlSrv{
		path:     path,
		dist:     dist,
		cli_path: cli_path,
		fs:       fs,
		backend:  backends.New(fs),
		provider: providers.New("lxc"),
	}
	return strm
}


// ## LXC IMPLEMENTATION ## //
func (c *LdlSrv) Create(name string) map[string]string {
	return c.provider.Create(name, c.fs, c.dist)
}

func (c *LdlSrv) Start(name string) map[string]string {
	//todo: maybe set static ip by using lxc.network.ipv4 = 10.0.3.211 on cfg
	return c.provider.Start(name)
}

func (c *LdlSrv) Stop(name string) map[string]string {
	return c.provider.Stop(name)
}

func (c *LdlSrv) Attach(name string) {
	c.provider.Attach(name)
}

func (c *LdlSrv) List() map[string]string {
	return c.provider.List()
}

func (c *LdlSrv) Freeze(name string) map[string]string {
	return c.provider.Freeze(name)
}

func (c *LdlSrv) Unfreeze(name string) map[string]string {
	return c.provider.Unfreeze(name)
}

func (c *LdlSrv) Info(name string) map[string]string {
	return c.provider.Info(name)
}

func (c *LdlSrv) Clone(from string, to string) map[string]string {
	c.Stop(from)
	result := c.provider.Clone(from, to)
	c.Start(from)
	return result
}

func (c *LdlSrv) Destroy(name string) map[string]string {
	c.Stop(name)
	if res := c.backend.Destroy(name); res["status"] != "ok" {
		return res
	}
	return c.provider.Destroy(name)
}

func (c *LdlSrv) Exec(name string, cmd string) map[string]string {
	return c.provider.Exec(name, cmd)
}


// ## SERVER IMPLEMENTATION ## //


func (c *LdlSrv) Commit(name string) map[string]string {
	c.Stop(name)
	result := c.backend.Snapshot(name, c.getSnapNum(name))
	c.Start(name)
	return result
}

func (c *LdlSrv) Log(name string) map[string]string {
	return c.backend.Snapshots(name)
}

func (c *LdlSrv) Push(name string) map[string]string {
	snapshots := c.Log(name)
	if snapshots["status"] != "ok" {
		return snapshots
	}

	snap_list := strings.Split(snapshots["message"], "\n")
	repo_path := fmt.Sprintf("%s/%s", c.path, name)
	snap_prev := ""
	snap_num := ""
	filename := ""

	os.MkdirAll(repo_path, 0755)

	for _, snap_name := range snap_list {
		snap_num = strings.Replace(snap_name, "snap", "", -1)
		filename = fmt.Sprintf("%s/%s.img", repo_path, snap_num)
		//fmt.Println(">", filename)
		if helpers.FileExists(filename) {
			snap_prev = snap_name
			continue
		}

		fmt.Println(">>>", snap_name)

		if snap_name == "snap0" {
			if dump := c.backend.DumpFull(name, filename); dump["status"] != "ok" {
				return dump
			}
		} else {
			if dump := c.backend.DumpIncr(name, snap_prev, name, snap_name, filename); dump["status"] != "ok" {
				return dump
			}
		}
		snap_prev = snap_name
	}

	return map[string]string{"status": "ok", "message": "success"}
}

func (c *LdlSrv) Export(name string, ssh string) {
	rsync, _ := exec.LookPath("rsync")
	src := fmt.Sprintf("%s/%s/", c.path, name)
	dst := fmt.Sprintf("%s:%s/%s", ssh, c.cli_path, name)
	syscall.Exec(rsync, []string{"rsync", "-auv", src, dst}, os.Environ())
}

func (c *LdlSrv) Share() {
	fmt.Println("Images shared at 0.0.0.0:8182")
	err := http.ListenAndServe(":8182", http.FileServer(http.Dir(c.path)))
	if err != nil {
		log.Errorf("unable to serve: %s", err)
	}
}

func (c *LdlSrv) getSnapNum(name string) int {
	log := c.Log(name)
	msg := log["message"]
	if log["status"] == "ok" && msg != "" && msg != "no datasets available" {
		return len(strings.Split(log["message"], "\n"))
	}
	return 0
}
