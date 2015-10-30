## Installation on Slave server

### Requirements

* Ubuntu >= 14.04
* LXC >= 1.0.7
* BTRFS >= 3.12


### Installation

	apt-get update && apt-get upgrade -y
	apt-get install -y lxc btrfs-tools iptables-persistent fail2ban

	modprobe btrfs
	mkfs.btrfs -f /dev/sdb
	echo "/dev/sdb /var/lib/lxc btrfs defaults 0 0" >> /etc/fstab
	mount /var/lib/lxc
	btrfs subvolume list /var/lib/lxc


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
