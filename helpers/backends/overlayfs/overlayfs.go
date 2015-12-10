package overlayfs

import (
	"fmt"

	"github.com/LPgenerator/Ldld/helpers"
)


var (
	// Server commands
	OVERLAYFS_SNAPSHOT = `lxc-snapshot -n %s 2>/dev/null`
	OVERLAYFS_DUMP_FULL = `tar --numeric-owner -cPf %s /var/lib/lxcsnaps/%s/snap0/ 2>/dev/null`
	OVERLAYFS_DUMP_INCR = `tar --numeric-owner -cPf %s /var/lib/lxcsnaps/%s/%s 2>/dev/null`
	OVERLAYFS_GET_SNAPSHOTS = `find /var/lib/lxcsnaps/%s/ -type d -name 'snap*' 2>/dev/null | xargs -I 1 basename 1 |sort -t p -k 3 -g | cut -d'-' -f2`
	OVERLAYFS_DEL_SNAPSHOTS = `rm -rf /var/lib/lxcsnaps/%s`

	// Client commands
	OVERLAYFS_CLONE_FS_FROM = `find /var/lib/lxcsnaps/%s -type d -name '%s' 2>/dev/null | sort -t p -k 3 -g | tail -1`
	OVERLAYFS_CLONE_FS = `rsync -au --exclude "config" --exclude "rootfs" %s/ /var/lib/lxc/%s`
	OVERLAYFS_SNAP_EXISTS = `ls -d /var/lib/lxcsnaps/%s/snap%d 2>/dev/null`
	OVERLAYFS_IMPORT = `tar xpf %s/%s/%d.img -C / 2>/dev/null`
)


type Overlayfs struct {
}


//
// SERVER
//
func (o Overlayfs) Snapshot(ct string, num int) map[string]string {
	snaps := o.Snapshots(ct)
	res := helpers.ExecRes(OVERLAYFS_SNAPSHOT, ct)
	if snaps["message"] == "" {
		helpers.ExecRes(`mount -o bind /var/lib/lxc/%s/rootfs /var/lib/lxcsnaps/%s/snap0/rootfs`, ct, ct)
	}
	return res
}



func (o Overlayfs) DumpFull(ct string, dst string) map[string]string {
	return helpers.ExecRes(OVERLAYFS_DUMP_FULL, dst, ct)
}


func (o Overlayfs) DumpIncr(ct1 string, prev string, ct2 string, curr string, dst string) map[string]string {
	return helpers.ExecRes(OVERLAYFS_DUMP_INCR, dst, ct1, curr)
}


func (o Overlayfs) Snapshots(ct string) map[string]string {
	return helpers.ExecRes(OVERLAYFS_GET_SNAPSHOTS, ct)
}


func (o Overlayfs) Destroy(ct string) map[string]string {
	return helpers.ExecRes(OVERLAYFS_DEL_SNAPSHOTS, ct)
}


//
// CLIENT
//
func (o Overlayfs) Optimize(ct string) map[string]string {
	// not implemented
	return map[string]string{"status": "ok", "message": "success"}
}


func (o Overlayfs) Mount(ct string) map[string]string {
	// todo: not implemented
	return map[string]string{"status": "ok", "message": "success"}
}


func (o Overlayfs) GetSnapshotByTemplate(template string, number string) map[string]string {
	if number == "" { number = "snap*" }
	return helpers.ExecRes(OVERLAYFS_CLONE_FS_FROM, template, number)
}


func (o Overlayfs) Clone(from string, to string) map[string]string {
	return helpers.ExecRes(OVERLAYFS_CLONE_FS, from, to)
}


func (o Overlayfs) SnapshotIsExists(ct string, num int) map[string]string {
	return helpers.ExecRes(OVERLAYFS_SNAP_EXISTS, ct, num)
}


func (o Overlayfs) ImportImage(path string, dist string, num int) map[string]string {
	return helpers.ExecRes(OVERLAYFS_IMPORT, path, dist, num)
}


func (o Overlayfs) AfterCreate(template string, name string) map[string]string {
	//todo: check latest delta. it can be delta1 or delta2

	fs := fmt.Sprintf("overlayfs:/var/lib/lxcsnaps/%s/snap0/rootfs:/var/lib/lxc/%s/delta0", template, name)
	if !helpers.SaveLXCDirective(name, "lxc.rootfs", fs) {
		return map[string]string{"status": "error", "message": "update config"}
	}

	return helpers.SaveHostInfo(name, fmt.Sprintf("/var/lib/lxc/%s/delta0/etc", name))
}
