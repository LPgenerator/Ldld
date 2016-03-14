ARGS=$(filter-out $@,$(MAKECMDGOALS))
BRANCH=`git rev-parse --abbrev-ref HEAD`
ENV=`basename "$PWD"`
NAME ?= ldld
PACKAGE_NAME ?= $(NAME)
PACKAGE_CONFLICT ?= $(PACKAGE_NAME)-beta
REVISION := $(shell git rev-parse --short HEAD || echo unknown)
LAST_TAG := $(shell git describe --tags --abbrev=0)
COMGPLv3S := $(shell echo `git log --oneline $(LAST_TAG)..HEAD | wc -l`)
VERSION := $(shell (cat VERSION || echo dev) | sed -e 's/^v//g')
ifneq ($(RELEASE),true)
    VERSION := $(shell echo $(VERSION)~beta.$(COMGPLv3S).g$(REVISION))
endif
ITTERATION := $(shell date +%s)
BUILD_PLATFORMS ?= -os="linux"
DEB_PLATFORMS ?= debian/wheezy debian/jessie ubuntu/precise ubuntu/trusty ubuntu/utopic ubuntu/vivid
DEB_ARCHS ?= amd64 i386 arm armhf
RPM_PLATFORMS ?= el/6 el/7 ol/6 ol/7
RPM_ARCHS ?= x86_64 i686 arm armhf

all: deps test lint toolchain build

deploy: build-and-deploy
	@rsync -auv out/deb/ldld_amd64.deb root@10.10.10.105:/tmp/$(PACKAGE_NAME)_$(PACKAGE_ARCH)-$(VERSION).deb
	@ssh root@10.10.10.105 "aptly repo add lpg /tmp/$(PACKAGE_NAME)_$(PACKAGE_ARCH)-$(VERSION).deb"
	@ssh root@10.10.10.105 "aptly publish update lpg"

run:
	dogo

help:
	# make run => run development server
	# make run-http-test => run wrk benchmarking
	# make register-fake-http => register backends
	# make session => run development session
	# make pull => pull updates from repo
	# make push => push changes to repo
	# make version - show information about current version
	# make deps - install all dependencies
	# make test - run project tests
	# make lint - check project code style
	# make toolchain - install crossplatform toolchain
	# make build - build project for all supported OSes
	# make package - package project using FPM

version: FORCE
	@echo Current version: $(VERSION)
	@echo Current iteration: $(ITTERATION)
	@echo Current revision: $(REVISION)

deps:
	# Installing dependencies...
	go get github.com/tools/godep
	go get -u github.com/golang/lint/golint
	go get github.com/mitchellh/gox
	go get golang.org/x/tools/cmd/cover
	godep restore

toolchain:
	# Building toolchain...
	gox $(BUILD_PLATFORMS)

build:
	gox $(BUILD_PLATFORMS) \
		-ldflags "-X main.NAME $(PACKAGE_NAME) -X main.VERSION $(VERSION) -X main.REVISION $(REVISION)" \
		-output="out/binaries/$(NAME)-{{.OS}}-{{.Arch}}"

lint:
	# Checking project code style...
	golint ./... | grep -v "be unexported"

test:
	# Running tests...
	go test ./... -cover

build-and-deploy:
	make build BUILD_PLATFORMS="-os=linux -arch=amd64"
	make package-deb-fpm ARCH=amd64 PACKAGE_ARCH=amd64
	make package-rpm-fpm ARCH=amd64 PACKAGE_ARCH=amd64

package: package-deps package-deb package-rpm

package-deb:
	# Building Debian compatible packages...
	make package-deb-fpm ARCH=amd64 PACKAGE_ARCH=amd64
	make package-deb-fpm ARCH=386 PACKAGE_ARCH=i386
	make package-deb-fpm ARCH=arm PACKAGE_ARCH=arm
	make package-deb-fpm ARCH=arm PACKAGE_ARCH=armhf

package-rpm:
	# Building RedHat compatible packages...
	make package-rpm-fpm ARCH=amd64 PACKAGE_ARCH=amd64
	make package-rpm-fpm ARCH=386 PACKAGE_ARCH=i686
	make package-rpm-fpm ARCH=arm PACKAGE_ARCH=arm
	make package-rpm-fpm ARCH=arm PACKAGE_ARCH=armhf

package-deps:
	# Installing packaging dependencies...
	gem install fpm

package-deb-fpm:
	@mkdir -p out/deb/
	fpm -s dir -t deb -n $(PACKAGE_NAME) -v $(VERSION) \
		-p out/deb/$(PACKAGE_NAME)_$(PACKAGE_ARCH).deb \
		--deb-priority optional --category admin \
		--force \
		--deb-compression bzip2 \
		--url https://github.com/LPgenerator/Ldld \
		--description "Ldld - Ldld is a simple and open source tool for shipping and running distributed containers." \
		-m "GoTLiuM InSPiRiT <gotlium@gmail.com>" \
		--license "GPLv3" \
		--vendor "github.com/gotlium" \
		--conflicts $(PACKAGE_CONFLICT) \
		--provides ldld \
		--replaces ldld \
		--after-install packaging/root/usr/share/ldld/post-install \
		--before-remove packaging/root/usr/share/ldld/post-install \
		-a $(PACKAGE_ARCH) \
		packaging/root/=/ \
		out/binaries/$(NAME)-linux-$(ARCH)=/usr/bin/ldld

package-rpm-fpm:
	@mkdir -p out/rpm/
	fpm -s dir -t rpm -n $(PACKAGE_NAME) -v $(VERSION) \
		-p out/rpm/$(PACKAGE_NAME)_$(PACKAGE_ARCH).rpm \
		--rpm-compression bzip2 --rpm-os linux \
		--force \
		--url https://github.com/LPgenerator/Ldld \
		--description "Ldld - Ldld is a simple and open source tool for shipping and running distributed containers." \
		-m "GoTLiuM InSPiRiT <gotlium@gmail.com>" \
		--license "GPLv3" \
		--vendor "github.com/gotlium" \
		--conflicts $(PACKAGE_CONFLICT) \
		--provides ldld \
		--replaces ldld \
		-a $(PACKAGE_ARCH) \
		packaging/root/=/ \
		out/binaries/$(NAME)-linux-$(ARCH)=/usr/bin/ldld

pull:
	@git pull origin `git rev-parse --abbrev-ref HEAD`
	@git log --name-only -1|grep migrations >& /dev/null && ./manage.py migrate --noinput || true
	@test -f touch.reload && touch touch.reload || true

push:
	@git status --porcelain|grep -v '??' && (echo '\033[0;32mCommit message:\033[0m' && MSG=`rlwrap -o -S "> " cat` && git commit -am "$$MSG") || true
	@git push origin $(BRANCH) || (git pull origin $(BRANCH) && git push origin $(BRANCH))

FORCE:
