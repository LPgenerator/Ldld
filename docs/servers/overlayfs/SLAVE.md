## Installation on Slave server (Not implemented)

### Requirements

* Ubuntu >= 14.04
* LXC >= 1.0.7
* OverlayFS >= 3.13-20140303


### Installation

	apt-get update && apt-get upgrade -y
	apt-get install -y lxc iptables-persistent fail2ban

	# download image which are you are using for CT
	lxc-create -t ubuntu -n base
	modprobe overlayfs


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
