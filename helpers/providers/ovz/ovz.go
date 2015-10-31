package ovz

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/LPgenerator/Ldld/helpers"
)

type Ovz struct {
}

var (
	CREATE_CT = `lxc-create -B %s -t %s -n %s`
	START_CT = `vzctl start %s`
	STOP_CT = `vzctl stop %s`
	LIST_CT = `vzlist -o ctid,status,ip,onboot,physpages,swappages`
	FREEZE_CT = `vzctl exec %s --suspend`
	UNFREEZE_CT = `vzctl exec %s --resume`
	INFO_CT = `lxc-info -n %s`
	RM_CT = `vzctl destroy %s`
	EXEC_CT = `vzctl exec %s %s`
	CLONE_CT = `lxc-clone -s %s %s`
)

func (o Ovz) Create(name string, fs string, dist string) map[string]string {
	return helpers.ExecRes(CREATE_CT, fs, dist, name)
}

func (o Ovz) Start(name string) map[string]string {
	return helpers.ExecRes(START_CT, name)
}

func (o Ovz) Stop(name string) map[string]string {
	return helpers.ExecRes(STOP_CT, name)
}

func (o Ovz) Attach(name string) {
	lxc_attach, _ := exec.LookPath("vzctl")
	syscall.Exec(lxc_attach, []string{"vzctl", "enter", name}, os.Environ())
}

func (o Ovz) List() map[string]string {
	return helpers.ExecRes(LIST_CT)
}

func (o Ovz) Freeze(name string) map[string]string {
	return helpers.ExecRes(FREEZE_CT, name)
}

func (o Ovz) Unfreeze(name string) map[string]string {
	return helpers.ExecRes(UNFREEZE_CT, name)
}

func (o Ovz) Info(name string) map[string]string {
	return helpers.ExecRes(INFO_CT, name)
}

func (o Ovz) Clone(src string, dst string) map[string]string {
	return helpers.ExecRes(CLONE_CT, src, dst)
}

func (o Ovz) Destroy(name string) map[string]string {
	return helpers.ExecRes(RM_CT, name)

}

func (o Ovz) Exec(name string, cmd string) map[string]string {
	return helpers.ExecRes(EXEC_CT, name, cmd)
}
