package btrfs

import (
	"github.com/LPgenerator/Ldld/helpers"
)


var (
	// Server commands
	BTRFS_SNAPSHOT = `btrfs subvolume snapshot -r /var/lib/lxc/%s/rootfs /var/lib/lxc/%s/rootfs-snap%d`
	BTRFS_DUMP_FULL = `btrfs send /var/lib/lxc/%s/rootfs-snap0 > %s`
	BTRFS_DUMP_INCR = `btrfs send -p /var/lib/lxc/%s/rootfs-%s /var/lib/lxc/%s/rootfs-%s > %s`
	BTRFS_GET_SNAPSHOTS = `find /var/lib/lxc/%s/ -type d -name 'rootfs-snap*' | xargs -I 1 basename 1 | cut -d'-' -f2`
	BTRFS_DEL_SNAPSHOTS = `find /var/lib/lxc/%s/ -type d -name 'rootfs-snap*' | xargs -I 1 btrfs subvolume delete 1`

	// Client commands
	BTRFS_CLONE_FS_FROM = `find /var/lib/lxc/%s/ -type d -name rootfs-%s|tail -1`
	BTRFS_CLONE_FS = `btrfs subvolume snapshot %s /var/lib/lxc/%s/rootfs`
	BTRFS_SNAP_EXISTS = `ls -d /var/lib/lxc/%s/rootfs-snap%d/`
	BTRFS_IMPORT = `cat %s/%s/%d.img | btrfs receive /var/lib/lxc/%s`
)


type Btrfs struct {
}


//
// SERVER
//
func (b Btrfs) Snapshot(ct string, num int) map[string]string {
	return helpers.ExecRes(BTRFS_SNAPSHOT, ct, ct, num)
}


func (b Btrfs) DumpFull(ct string, dst string) map[string]string {
	return helpers.ExecRes(BTRFS_DUMP_FULL, ct, dst)
}


func (b Btrfs) DumpIncr(ct1 string, prev string, ct2 string, curr string, dst string) map[string]string {
	return helpers.ExecRes(BTRFS_DUMP_INCR, ct1, prev, ct2, curr, dst)
}


func (b Btrfs) Snapshots(ct string) map[string]string {
	return helpers.ExecRes(BTRFS_GET_SNAPSHOTS, ct)
}


func (b Btrfs) Destroy(ct string) map[string]string {
	return helpers.ExecRes(BTRFS_DEL_SNAPSHOTS, ct)
}


//
// CLIENT
//
func (b Btrfs) Optimize(ct string) map[string]string {
	// not implemented
	return map[string]string{"status": "ok", "message": "success"}
}


func (b Btrfs) Mount(ct string) map[string]string {
	// todo: not implemented
	return map[string]string{"status": "ok", "message": "success"}
}


func (b Btrfs) GetSnapshotByTemplate(template string, number string) map[string]string {
	if number == "" { number = "snap*" }
	return helpers.ExecRes(BTRFS_CLONE_FS_FROM, template, number)
}


func (b Btrfs) Clone(from string, to string) map[string]string {
	return helpers.ExecRes(BTRFS_CLONE_FS, from, to)
}


func (b Btrfs) SnapshotIsExists(ct string, num int) map[string]string {
	return helpers.ExecRes(BTRFS_SNAP_EXISTS, ct, num)
}


func (b Btrfs) ImportImage(path string, dist string, num int) map[string]string {
	helpers.ExecRes("mkdir -p /var/lib/lxc/%s", dist)
	return helpers.ExecRes(BTRFS_IMPORT, path, dist, num, dist)
}


func (b Btrfs) AfterCreate(template string, name string) map[string]string {
	// not implemented
	return map[string]string{"status": "ok", "message": "success"}
}
