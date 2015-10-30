## Installation on Primary server

### Requirements

* Ubuntu >= 14.04
* LXC >= 1.0.7
* ZFS >= 0.6.5

_Note: use zfs in production only with mirror_


### Installation

	apt-get update
	apt-get upgrade
	apt-get install -y python-software-properties software-properties-common
	apt-add-repository -y ppa:zfs-native/stable
	apt-get update
	apt-get install -y lxc ubuntu-zfs nginx fail2ban
	apt-get clean
	
	# nginx configuration
	rm /usr/share/nginx/html/*.html
	cat > /etc/nginx/sites-enabled/default << END
	server {
		root /usr/share/nginx/html;
		listen 80 default_server;
		server_name _;
	
		location / {
			if (\$http_user_agent !~ (Go|Wget) ) {
				return 403;
			}
			autoindex on;
		}
	}
	END
	cat > /etc/lxc/lxc.conf << END
	lxc.lxcpath = /var/lib/lxc
	lxc.bdev.zfs.root = ldld/lxc
	END
	
	service nginx restart

	# fs configuration
	modprobe zfs
	zpool create ldld -f /dev/sdb
	zfs create ldld/lxc
	zfs set dedup=off ldld
	zfs set compression=off ldld
	zpool set listsnapshots=off ldld


### Create CT

	ldld create web


### Start CT

	ldld start web
	ldld info web


### Attach to CT

	ldld attach web


### Commit and Push to Repo

	ldld commit web
	ldld push web
