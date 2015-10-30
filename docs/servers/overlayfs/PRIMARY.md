## Installation on Primary server (Not implemented)

### Requirements

* Ubuntu >= 14.04
* LXC >= 1.0.7
* OverlayFS >= 3.13-20140303


### Installation

	apt-get update
	apt-get upgrade
	apt-get install -y lxc nginx fail2ban
	apt-get clean
	modprobe overlayfs
	
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
	lxc.bdev.lvm.vg = vg0
	END

	service nginx restart


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
