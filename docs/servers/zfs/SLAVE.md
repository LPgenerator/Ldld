## Installation on Slave server

### Requirements

* Ubuntu >= 14.04
* LXC >= 1.0.7
* ZFS >= 0.6.5


### Installation

	apt-get update
	apt-get install -y python-software-properties software-properties-common
	apt-add-repository ppa:zfs-native/stable
	apt-get update
	apt-get install lxc ubuntu-zfs iptables-persistent fail2ban

	cat > /etc/lxc/lxc.conf << END
	lxc.lxcpath = /var/lib/lxc
	lxc.bdev.zfs.root = ldld/lxc
	END

	modprobe zfs
	zpool create ldld -f /dev/sdb
	zfs create ldld/lxc
	zfs set dedup=off ldld
	zfs set compression=off ldld
	zpool set listsnapshots=off ldld


### Download Images

	ldld pull web
	ldld pull celery


### Create CT from Image

	ldld create web-1 web
	ldld create web-2 web:1


### Allow autostart on boot

	ldld autostart web-1 1
	ldld autostart web-2 1


### Start CT

	ldld start web-1
	ldld start web-2


### Info about CT

	ldld info web-1


### Attach to CT

	ldld attach web-1
