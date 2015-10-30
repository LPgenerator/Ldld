package zfs

import "github.com/LPgenerator/Ldld/helpers"


var (
	// Server commands
	ZFS_DEL_ROOTFS = `zfs destroy -R ldld/lxc/%s`
	ZFS_SNAPSHOT = `zfs snapshot ldld/lxc/%s@snap%d`
	ZFS_DUMP_FULL = `zfs send ldld/lxc/%s@snap0 > %s`
	ZFS_DUMP_INCR = `zfs send -i ldld/lxc/%s@%s ldld/lxc/%s@%s > %s`
	ZFS_GET_SNAPSHOTS = `zfs list -t snapshot|grep ldld/lxc/%s@snap|awk '{print $1}'| cut -d'@' -f2` //| tac
	ZFS_DEL_SNAPSHOTS = `zfs list -t snapshot|grep ldld/lxc/%s@snap|awk '{print $1}'| xargs -I 1 zfs destroy -R 1`

	// Client commands
	ZFS_OPTIMIZE_FS_SYNC = `zfs set sync=disabled ldld/lxc/%s`
	ZFS_OPTIMIZE_FS_CKSUM = `zfs set checksum=off ldld/lxc/%s`
	ZFS_OPTIMIZE_FS_ATIME = `zfs set atime=off ldld/lxc/%s`
	ZFS_SET_FS_MOUNT_POINT = `zfs set mountpoint=/var/lib/lxc/%s/rootfs ldld/lxc/%s`
	ZFS_CLONE_FS_FROM = `zfs list -t snapshot | grep -P '%s@%s$' | tail -1 | awk "{print \$1}"`
	ZFS_CLONE_FS = `zfs clone %s ldld/lxc/%s`
	ZFS_SNAP_EXISTS = `zfs list -t snapshot|grep ldld/lxc/%s@snap%d`
	ZFS_IMPORT = `cat %s/%s/%d.img | zfs receive ldld/lxc/%s`
)


type Zfs struct {
}


//
// SERVER
//
func (z Zfs) Snapshot(ct string, num int) map[string]string {
	return helpers.ExecRes(ZFS_SNAPSHOT, ct, num)
}


func (z Zfs) DumpFull(ct string, dst string) map[string]string {
	return helpers.ExecRes(ZFS_DUMP_FULL, ct, dst)
}


func (z Zfs) DumpIncr(ct1 string, prev string, ct2 string, curr string, dst string) map[string]string {
	return helpers.ExecRes(ZFS_DUMP_INCR, ct1, prev, ct2, curr, dst)
}


func (z Zfs) Snapshots(ct string) map[string]string {
	return helpers.ExecRes(ZFS_GET_SNAPSHOTS, ct)
}


func (z Zfs) Destroy(ct string) map[string]string {
	if res := helpers.ExecRes(ZFS_DEL_SNAPSHOTS, ct); res["status"] != "ok" {
		return res
	}
	if res := helpers.ExecRes(ZFS_DEL_ROOTFS, ct); res["status"] != "ok" {
		return res
	}
	return map[string]string{"status": "ok", "message": "success"}
}


//
// CLIENT
//
func (z Zfs) Optimize(ct string) map[string]string {
	if res := helpers.ExecRes(ZFS_OPTIMIZE_FS_SYNC, ct); res["status"] != "ok" {
		return res
	}
	if res := helpers.ExecRes(ZFS_OPTIMIZE_FS_CKSUM, ct); res["status"] != "ok" {
		return res
	}
	if res := helpers.ExecRes(ZFS_OPTIMIZE_FS_ATIME, ct); res["status"] != "ok" {
		return res
	}
	return map[string]string{"status": "ok", "message": "success"}
}


func (z Zfs) Mount(ct string) map[string]string {
	return helpers.ExecRes(ZFS_SET_FS_MOUNT_POINT, ct, ct)
}


func (z Zfs) GetSnapshotByTemplate(template string, number string) map[string]string {
	if number == "" { number = ".*" }
	return helpers.ExecRes(ZFS_CLONE_FS_FROM, template, number)
}


func (z Zfs) Clone(from string, to string) map[string]string {
	return helpers.ExecRes(ZFS_CLONE_FS, from, to)
}


func (z Zfs) SnapshotIsExists(ct string, num int) map[string]string {
	return helpers.ExecRes(ZFS_SNAP_EXISTS, ct, num)
}


func (z Zfs) ImportImage(path string, dist string, num int) map[string]string {
	return helpers.ExecRes(ZFS_IMPORT, path, dist, num, dist)
}


func (z Zfs) AfterCreate(template string, name string) map[string]string {
	// not implemented
	return map[string]string{"status": "ok", "message": "success"}
}
