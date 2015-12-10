package lvm

import (
	"github.com/LPgenerator/Ldld/helpers"
)

// todo: incremental backups, restore from images, lv size & etc


var (
	// Server commands
	LVM_SNAPSHOT = `lvcreate -s -n %s-snap%d -l 1 /dev/vg0/%s`
	LVM_DUMP_FULL = `cat /dev/vg0/%s-snap0 | gzip > %s`
	LVM_DUMP_INCR = `cat /dev/vg0/%s-%s | gzip > %s`
	LVM_GET_SNAPSHOTS = `lvs vg0| grep %s-snap | awk '{print $1}' | cut -d'-' -f2`
	LVM_DEL_SNAPSHOTS = `lvs vg0| grep %s-snap | awk '{print $1}' | xargs -I 1 lvremove -f vg0/1`

	// Client commands
	LVM_CLONE_FS_FROM = `lvs vg0| grep %s-snap| awk '{print $1}' |sort -t p -k 2 -g | grep -P 'snap%s$'`
	LVM_CLONE_FS = `?` // how to clone without new fs size
	LVM_SNAP_EXISTS = `lvs vg0| grep %s-snap%d`
	LVM_IMPORT = `cat %s/%s/%d.img | gunzip > /dev/vg0/%s-%s` // create lv at first with unknown fs size
)


type Lvm struct {
}


//
// SERVER
//
func (b Lvm) Snapshot(ct string, num int) map[string]string {
	return helpers.ExecRes(LVM_SNAPSHOT, ct, num, ct)
}


func (b Lvm) DumpFull(ct string, dst string) map[string]string {
	return helpers.ExecRes(LVM_DUMP_FULL, ct, dst)
}


func (b Lvm) DumpIncr(ct1 string, prev string, ct2 string, curr string, dst string) map[string]string {
	return helpers.ExecRes(LVM_DUMP_INCR, ct1, curr, dst)
}


func (b Lvm) Snapshots(ct string) map[string]string {
	return helpers.ExecRes(LVM_GET_SNAPSHOTS, ct)
}


func (b Lvm) Destroy(ct string) map[string]string {
	return helpers.ExecRes(LVM_DEL_SNAPSHOTS, ct)
}


//
// CLIENT
//
func (b Lvm) Optimize(ct string) map[string]string {
	// not implemented
	return map[string]string{"status": "ok", "message": "success"}
}


func (b Lvm) Mount(ct string) map[string]string {
	// not implemented
	return map[string]string{"status": "ok", "message": "success"}
}


func (b Lvm) GetSnapshotByTemplate(template string, number string) map[string]string {
	if number == "" { number = ".*" }
	return helpers.ExecRes(LVM_CLONE_FS_FROM, template, number)
}


func (b Lvm) Clone(from string, to string) map[string]string {
	return helpers.ExecRes(LVM_CLONE_FS, from, to)
}


func (b Lvm) SnapshotIsExists(ct string, num int) map[string]string {
	return helpers.ExecRes(LVM_SNAP_EXISTS, ct, num)
}


func (b Lvm) ImportImage(path string, dist string, num int) map[string]string {
	helpers.ExecRes("mkdir -p /var/lib/lxc/%s", dist)
	return helpers.ExecRes(LVM_IMPORT, path, dist, num, dist)
}
