package backends

import (
	"github.com/LPgenerator/Ldld/helpers/backends/zfs"
	"github.com/LPgenerator/Ldld/helpers/backends/btrfs"
	//"github.com/LPgenerator/Ldld/helpers/backends/overlayfs"
	//"github.com/LPgenerator/Ldld/helpers/backends/lvm"
)


type Fs interface {
	Snapshot(ct string, num int) map[string]string
	DumpFull(ct string, dst string) map[string]string
	DumpIncr(ct1 string, prev string, ct2 string, curr string, dst string) map[string]string
	Snapshots(ct string) map[string]string
	Destroy(ct string) map[string]string
	Optimize(ct string) map[string]string
	Mount(ct string) map[string]string
	GetSnapshotByTemplate(template string, number string) map[string]string
	Clone(from string, to string) map[string]string
	SnapshotIsExists(ct string, num int) map[string]string
	ImportImage(path string, dist string, num int) map[string]string
	AfterCreate(template string, name string) map[string]string
}


func New(backend string) (Fs) {
	fs := []Fs{zfs.Zfs{}}
	if backend == "btrfs" {
		fs = []Fs{btrfs.Btrfs{}}
	}
	/*
	if backend == "overlayfs" {
		fs = []Fs{overlayfs.Overlayfs{}}
	}
	if backend == "lvm" {
		fs = []Fs{lvm.Lvm{}}
	}
	*/
	return fs[0]
}
