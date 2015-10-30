#!/bin/bash

cat > /etc/apt/sources.list << END
deb http://ubuntu.uz/ubuntu/ trusty main restricted
deb http://ubuntu.uz/ubuntu/ trusty-updates main restricted
deb http://ubuntu.uz/ubuntu/ trusty universe
deb http://ubuntu.uz/ubuntu/ trusty-updates universe
deb http://ubuntu.uz/ubuntu/ trusty multiverse
deb http://ubuntu.uz/ubuntu/ trusty-updates multiverse
deb http://ubuntu.uz/ubuntu/ trusty-backports main restricted universe multiverse
deb http://ubuntu.uz/ubuntu/ trusty-security main restricted
deb http://ubuntu.uz/ubuntu/ trusty-security universe
deb http://ubuntu.uz/ubuntu/ trusty-security multiverse

END

apt-get update
apt-get upgrade
apt-get install -y lxc nginx golang git

if [ "`hostname`" == "ldl-master" ]; then
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
lxc-fs = "overlayfs"
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
lxc-fs = "overlayfs"
END
fi

# MIRROR="http://ubuntu.uz/ubuntu" SECURITY_MIRROR="http://ubuntu.uz/ubuntu" lxc-create -t ubuntu -n test-1 -B zfs --zfsroot lpg/lxc
# lxc-destroy -f -n test-1
