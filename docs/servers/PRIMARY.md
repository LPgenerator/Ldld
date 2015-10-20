## Installation on Primary server

### Requirements

* Ubuntu >= 14.04
* LXC >= 1.0.7
* ZFS >= 0.6.5


### Installation

	apt-get update
	apt-get upgrade
	apt-get install -y python-software-properties software-properties-common
	apt-add-repository -y ppa:zfs-native/stable
	apt-get update
	apt-get install -y lxc ubuntu-zfs nginx fail2ban
	
	# nginx configuration
	rm /usr/share/nginx/html/*.html
	touch /etc/nginx/sites-enabled/default && cat > /etc/nginx/sites-enabled/default << END
	server {
		root /usr/share/nginx/html;
		listen 80 default_server;
		server_name _;
	
		location / {
			autoindex on;
		}
	}
	END
	service nginx restart

	# fs configuration
	modprobe zfs
	zpool create lpg /dev/sdb
	zfs create lpg/lxc


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
