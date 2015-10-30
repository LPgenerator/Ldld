## Installation on Primary server

### Requirements

* Ubuntu >= 14.04
* LXC >= 1.0.7
* BTRFS >= 3.12


### Installation

	apt-get update
	apt-get upgrade
	apt-get install -y lxc btrfs-tools nginx fail2ban
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

	service nginx restart

	# fs configuration
	modprobe btrfs
	mkfs.btrfs -f /dev/sdb
	echo "/dev/sdb /var/lib/lxc btrfs defaults 0 0" >> /etc/fstab
	mount /var/lib/lxc
	btrfs subvolume list /var/lib/lxc


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
