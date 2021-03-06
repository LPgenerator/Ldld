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
)

var (
	CREATE_CT = `lxc-create -B %s -t %s -n %s`
	START_CT = `lxc-start -d -n %s`
	EXEC_CT = `lxc-attach -n %s -- %s`
	STOP_CT = `lxc-stop -n %s`
	LIST_CT = `lxc-ls -f -F name,state,ipv4,ipv6,autostart,pid,ram,swap`
	CLONE_CT = `lxc-clone -s %s %s`
	FREEZE_CT = `lxc-freeze -n %s`
	UNFREEZE_CT = `lxc-unfreeze -n %s`
	INFO_CT = `lxc-info -n %s`
	RM_CT = `lxc-destroy -f -n %s`
)


type LdlSrv struct{
	path      string
	dist      string
	cli_path  string
	fs        string
	backend   backends.Fs
}

func New(path string, dist string, cli_path string, fs string) (*LdlSrv) {
	strm := &LdlSrv{
		path:     path,
		dist:     dist,
		cli_path: cli_path,
		fs:       fs,
		backend:  backends.New(fs),
	}
	return strm
}

func (c *LdlSrv) getSnapNum(name string) int {
	log := c.Log(name)
	msg := log["message"]
	if log["status"] == "ok" && msg != "" && msg != "no datasets available" {
		return len(strings.Split(log["message"], "\n"))
	}
	return 0
}


// ## LXC IMPLEMENTATION ## //
func (c *LdlSrv) Create(name string) map[string]string {
	return helpers.ExecRes(CREATE_CT, c.fs, c.dist, name)
}

func (c *LdlSrv) Start(name string) map[string]string {
	//todo: maybe set static ip by using lxc.network.ipv4 = 10.0.3.211 on cfg
	return helpers.ExecRes(START_CT, name)
}

func (c *LdlSrv) Stop(name string) map[string]string {
	return helpers.ExecRes(STOP_CT, name)
}

func (c *LdlSrv) Attach(name string) {
	lxc_attach, _ := exec.LookPath("lxc-attach")
	syscall.Exec(lxc_attach, []string{"lxc-attach", "-n", name}, os.Environ())
}

func (c *LdlSrv) List() map[string]string {
	return helpers.ExecRes(LIST_CT)
}

func (c *LdlSrv) Freeze(name string) map[string]string {
	return helpers.ExecRes(FREEZE_CT, name)
}

func (c *LdlSrv) Unfreeze(name string) map[string]string {
	return helpers.ExecRes(UNFREEZE_CT, name)
}

func (c *LdlSrv) Info(name string) map[string]string {
	return helpers.ExecRes(INFO_CT, name)
}

func (c *LdlSrv) Commit(name string) map[string]string {
	c.Stop(name)
	result := c.backend.Snapshot(name, c.getSnapNum(name))
	c.Start(name)
	return result
}

func (c *LdlSrv) Log(name string) map[string]string {
	return c.backend.Snapshots(name)
}

func (c *LdlSrv) Clone(from string, to string) map[string]string {
	c.Stop(from)
	result := helpers.ExecRes(CLONE_CT, from, to)
	c.Start(from)
	return result
}

func (c *LdlSrv) Destroy(name string) map[string]string {
	c.Stop(name)
	if res := c.backend.Destroy(name); res["status"] != "ok" {
		return res
	}
	return helpers.ExecRes(RM_CT, name)
}

func (c *LdlSrv) Exec(name string, cmd string) map[string]string {
	return helpers.ExecRes(EXEC_CT, name, cmd)
}


// ## REPOSITORY SERVER IMPLEMENTATION ## //
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
