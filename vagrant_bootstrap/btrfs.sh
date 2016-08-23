#!/usr/bin/env bash

apt-get update
apt-get install -y lxc btrfs-tools golang git

if [ "`hostname`" == "ldl-master" ]; then
    apt-get install -y nginx
    rm /usr/share/nginx/html/*.html
    cat > /etc/nginx/sites-enabled/default << END
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
fi

# fs configuration
modprobe btrfs
mkfs.btrfs -f /dev/sdc || mkfs.btrfs -f /dev/sdb
echo "/dev/sdb /var/lib/lxc btrfs defaults 0 0" >> /etc/fstab
mount /var/lib/lxc

export GOPATH="/root/.go"
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

echo 'export GOPATH="/root/.go"' >> ~/.bashrc
echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' >> ~/.bashrc

go get github.com/tools/godep
go get github.com/mitchellh/gox

cd /vagrant
godep restore
mkdir -p /root/.go/src/gitlab.com/ayufan/
cd /root/.go/src/gitlab.com/ayufan/; git clone https://gitlab.com/ayufan/golang-cli-helpers.git; cd -
mkdir /root/.go/src/github.com/LPgenerator/
ln -sf /vagrant /root/.go/src/github.com/LPgenerator/Ldld
make build BUILD_PLATFORMS="-os=linux -arch=amd64"
ln -sf /vagrant/out/binaries/ldld-linux-amd64 /usr/local/bin/ldld

mkdir /etc/ldld/
if [ "`hostname`" == "ldl-master" ]; then
    cat > /etc/ldld/config.toml << END
api-address = "0.0.0.0:9090"
api-login = "ldl"
api-password = "7eNQ4iWLgDw4Q6w"
web-address = "0.0.0.0:9191"
web-login = "admin"
web-password = "7eNQ4iWLgDw4Q6w"
type = "server"
cli-repo-url = "http://48.44.44.44"
srv-path = "/usr/share/nginx/html"
cli-data-dir = "/usr/local/var/lib/ldl"
lxc-distro = "ubuntu"
END
else
    cat > /etc/ldld/config.toml << END
api-address = "0.0.0.0:9090"
api-login = "ldl"
api-password = "7eNQ4iWLgDw4Q6w"
web-address = "0.0.0.0:9191"
web-login = "admin"
web-password = "7eNQ4iWLgDw4Q6w"
type = "client"
cli-repo-url = "http://48.44.44.44"
srv-path = "/usr/share/nginx/html"
cli-data-dir = "/usr/local/var/lib/ldl"
lxc-distro = "ubuntu"
END
fi
