package overlayfs

import (
	"github.com/LPgenerator/Ldld/helpers"
)


var (
	// Server commands
	OVERLAYFS_SNAPSHOT = `lxc-snapshot -n %s`
	OVERLAYFS_DUMP_FULL = `tar hczf %s /var/lib/lxcsnaps/%s/snap0/`
	OVERLAYFS_DUMP_INCR = `tar czf %s /var/lib/lxcsnaps/%s/%s`
	OVERLAYFS_GET_SNAPSHOTS = `find /var/lib/lxcsnaps/%s/ -type d -name 'snap*' 2>/dev/null | xargs -I 1 basename 1 |sort -t p -k 2 -g | cut -d'-' -f2`
	OVERLAYFS_DEL_SNAPSHOTS = `find /var/lib/lxcsnaps/%s -type d -name 'snap*' | xargs -I 1 bash -c 'rm -rf 1'`

	// Client commands
	OVERLAYFS_CLONE_FS_FROM = `find /var/lib/lxcsnaps/%s -type d -name 'snap%s'|tail -1`
	OVERLAYFS_CLONE_FS = `lxc-clone -o %s -n %s -B overlayfs -s`
	OVERLAYFS_SNAP_EXISTS = `ls -d /var/lib/lxcsnaps/%s/snap%d`
	OVERLAYFS_IMPORT = `tar xzf %s/%s/%d.img -C /`
)


type Overlayfs struct {
}


//
// SERVER
//
func (b Overlayfs) Snapshot(ct string, num int) map[string]string {
	//snaps := b.Snapshots(ct)
	//res := helpers.ExecRes(OVERLAYFS_SNAPSHOT, ct)
	//if snaps["message"] == "" {
	//	helpers.ExecRes(`rmdir /var/lib/lxcsnaps/%s/snap0/rootfs; ln -sf /var/lib/lxc/%s/rootfs/ /var/lib/lxcsnaps/%s/snap0/rootfs`, ct, ct, ct)
	//}
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
	return helpers.ExecRes(OVERLAYFS_DEL_SNAPSHOTS, ct, ct)
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
	// todo: save to config
	return helpers.ExecRes(OVERLAYFS_CLONE_FS, from, to)
}


func (b Overlayfs) SnapshotIsExists(ct string, num int) map[string]string {
	return helpers.ExecRes(OVERLAYFS_SNAP_EXISTS, ct, num)
}


func (b Overlayfs) ImportImage(path string, dist string, num int) map[string]string {
	return helpers.ExecRes(OVERLAYFS_IMPORT, path, dist, num)
}
