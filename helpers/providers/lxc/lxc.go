package lxc

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/LPgenerator/Ldld/helpers"
)

type Lxc struct {
}

var (
	CREATE_CT = `lxc-create -B %s -t %s -n %s`
	START_CT = `lxc-start -d -n %s`
	STOP_CT = `lxc-stop -n %s` //+
	LIST_CT = `lxc-ls -f -F name,state,ipv4,autostart,memory,swap`
	FREEZE_CT = `lxc-freeze -n %s`
	UNFREEZE_CT = `lxc-unfreeze -n %s`
	INFO_CT = `lxc-info -n %s`
	RM_CT = `lxc-destroy -f -n %s`
	EXEC_CT = `lxc-attach -n %s -- %s`
	CLONE_CT = `lxc-clone -s %s %s`
)

func (l Lxc) Create(name string, fs string, dist string) map[string]string {
	return helpers.ExecRes(CREATE_CT, fs, dist, name)
}

func (l Lxc) Start(name string) map[string]string {
	return helpers.ExecRes(START_CT, name)
}

func (l Lxc) Stop(name string) map[string]string {
	return helpers.ExecRes(STOP_CT, name)
}

func (l Lxc) Attach(name string) {
	lxc_attach, _ := exec.LookPath("lxc-attach")
	syscall.Exec(lxc_attach, []string{"lxc-attach", "-n", name}, os.Environ())
}

func (l Lxc) List() map[string]string {
	return helpers.ExecRes(LIST_CT)
}

func (l Lxc) Freeze(name string) map[string]string {
	return helpers.ExecRes(FREEZE_CT, name)
}

func (l Lxc) Unfreeze(name string) map[string]string {
	return helpers.ExecRes(UNFREEZE_CT, name)
}

func (l Lxc) Info(name string) map[string]string {
	return helpers.ExecRes(INFO_CT, name)
}

func (l Lxc) Clone(src string, dst string) map[string]string {
	return helpers.ExecRes(CLONE_CT, src, dst)
}

func (l Lxc) Destroy(name string) map[string]string {
	return helpers.ExecRes(RM_CT, name)

}

func (l Lxc) Exec(name string, cmd string) map[string]string {
	return helpers.ExecRes(EXEC_CT, name, cmd)
}
