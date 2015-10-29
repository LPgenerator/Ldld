package backends


var (
	PORT_FORWARD = `iptables -t nat -A PREROUTING -p tcp --dport %s -j DNAT --to %s:%s`
	RM_PORT_FORWARD = `iptables -t nat -L PREROUTING -n -v --line-numbers| grep %s| awk '{print $1}'| xargs iptables -t nat -D PREROUTING`
	SET_FS_MOUNT_POINT = `zfs set mountpoint=/var/lib/lxc/%s/rootfs lpg/lxc/%s`
	OPTIMIZE_FS_SYNC = `zfs set sync=disabled lpg/lxc/%s`
	OPTIMIZE_FS_CKSUM = `zfs set checksum=off lpg/lxc/%s`
	OPTIMIZE_FS_ATIME = `zfs set atime=off lpg/lxc/%s`
	CLONE_FS = `zfs clone %s lpg/lxc/%s`
	CLONE_FS_FROM = `zfs list -t snapshot|grep %s@%s|tail -1|awk "{print \$1}"`
	WGET = `wget -c --retry-connrefused -t 0 %s/%s/%s -O %s/%s/%s`
	LinksRegexp = regexp.MustCompile(`">(.*?)/?</a>`)
	DESTROY_CT = `zfs destroy -rR lpg/lxc/%s >&/dev/null; lxc-destroy -f -n %s`
	EMPTY_CMD = ``
	MIGRATE_CFG = `scp /var/lib/lxc/%s/config %s:/var/lib/lxc/%s/`
	MIGRATE_ZFS = `zfs send lpg/lxc/%s@migrate | ssh %s zfs recv -F lpg/lxc/%s`
	MIGRATE_MP = `ssh %s zfs set mountpoint=/var/lib/lxc/%s/rootfs lpg/lxc/%s`
)


type Zfs struct {
}


func New() (*Zfs) {
	strm := &Zfs{
	}
	return strm
}

func (z *Zfs) Create() {
	//
}
//zfs_mount
//zfs_umount
//zfs_clone
//zfs_destroy
//zfs_create
