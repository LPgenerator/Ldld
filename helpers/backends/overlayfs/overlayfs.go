package overlayfs

import (
	"fmt"
	"github.com/LPgenerator/Ldld/helpers"
)


var (
	// Server commands
	OVERLAYFS_SNAPSHOT = `lxc-snapshot -n %s 2>/dev/null`
	OVERLAYFS_DUMP_FULL = `tar hczf --warning=no-file-changed --warning=no-file-removed -f %s /var/lib/lxcsnaps/%s/snap0/ 2>/dev/null`
	OVERLAYFS_DUMP_INCR = `tar czf %s /var/lib/lxcsnaps/%s/%s 2>/dev/null`
	OVERLAYFS_GET_SNAPSHOTS = `find /var/lib/lxcsnaps/%s/ -type d -name 'snap*' 2>/dev/null | xargs -I 1 basename 1 |sort -t p -k 3 -g | cut -d'-' -f2`
	OVERLAYFS_DEL_SNAPSHOTS = `rm -rf /var/lib/lxcsnaps/%s`

	// Client commands
	OVERLAYFS_CLONE_FS_FROM = `find /var/lib/lxcsnaps/%s -type d -name 'snap%s' 2>/dev/null | sort -t p -k 3 -g | tail -1`
	OVERLAYFS_CLONE_FS = `rsync -av --exclude "config" %s/ /var/lib/lxc/%s`
	OVERLAYFS_SNAP_EXISTS = `ls -d /var/lib/lxcsnaps/%s/snap%d 2>/dev/null`
	OVERLAYFS_IMPORT = `tar xzf %s/%s/%d.img -C / 2>/dev/null`
)


type Overlayfs struct {
}


//
// SERVER
//
func (b Overlayfs) Snapshot(ct string, num int) map[string]string {
	return helpers.ExecRes(OVERLAYFS_SNAPSHOT, ct)
}


func (b Overlayfs) DumpFull(ct string, dst string) map[string]string {
	return helpers.ExecRes(OVERLAYFS_DUMP_FULL, dst, ct)
}


func (b Overlayfs) DumpIncr(ct1 string, prev string, ct2 string, curr string, dst string) map[string]string {
	return helpers.ExecRes(OVERLAYFS_DUMP_INCR, dst, ct1, curr)
}


func (b Overlayfs) Snapshots(ct string) map[string]string {
	return helpers.ExecRes(OVERLAYFS_GET_SNAPSHOTS, ct)
}


func (b Overlayfs) Destroy(ct string) map[string]string {
	return helpers.ExecRes(OVERLAYFS_DEL_SNAPSHOTS, ct)
}


//
// CLIENT
//
func (b Overlayfs) Optimize(ct string) map[string]string {
	// not implemented
	return map[string]string{"status": "ok", "message": "success"}
}


func (b Overlayfs) Mount(ct string) map[string]string {
	// todo: not implemented
	return map[string]string{"status": "ok", "message": "success"}
}


func (b Overlayfs) GetSnapshotByTemplate(template string, number string) map[string]string {
	if number == "" { number = "*" }
	return helpers.ExecRes(OVERLAYFS_CLONE_FS_FROM, template, number)
}


func (b Overlayfs) Clone(from string, to string) map[string]string {
	return helpers.ExecRes(OVERLAYFS_CLONE_FS, from, to)
}


func (b Overlayfs) SnapshotIsExists(ct string, num int) map[string]string {
	return helpers.ExecRes(OVERLAYFS_SNAP_EXISTS, ct, num)
}


func (b Overlayfs) ImportImage(path string, dist string, num int) map[string]string {
	return helpers.ExecRes(OVERLAYFS_IMPORT, path, dist, num)
}


func (b Overlayfs) AfterCreate(name string) map[string]string {
	//todo: check latest delta. it can be delta1 or delta2
	//todo: rootfs now is hardcoded. need another way to do this
	//todo: set hostname & add to /etc/hosts
	//todo: use signals BeforeCreate and AfterCreate

	fs := fmt.Sprintf("overlayfs:/var/cache/lxc/trusty/rootfs-amd64:/var/lib/lxc/%s/delta0", name)
	if !helpers.SaveLXCDirective(name, "lxc.rootfs", fs) {
		return map[string]string{"status": "error", "message": "update config"}
	}
	return map[string]string{"status": "ok", "message": "success"}
}
