# Development environment

## 1. Install dependencies and Go runtime

### For Debian/Ubuntu
```bash
sudo apt-get install -y curl bison gcc
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source ~/.gvm/scripts/gvm
gvm install go1.4.3
gvm use go1.4.3
export GOROOT_BOOTSTRAP=$GOROOT
gvm install go1.5.1
gvm use go1.5.1
gvm pkgset create ldld
gvm pkgset use ldld
```

## 2. Download Ldld sources

```
go get github.com/LPgenerator/Ldld
cd ~/.gvm/pkgsets/go1.5.1/ldld/src/github.com/LPgenerator/Ldld
```

## 3. Install Ldld dependencies

This will download and restore all dependencies required to build `Ldld`:

```
make deps
```

## 4. Run Ldld

Normally you would use `Ldld`, in order to compile and run Go source use go toolchain:

```
go run main.go run
```

You can run `Ldld` in debug-mode:

```
go run --debug main.go run
```

## 5. Compile and install Ldld binary

```
go build
go install
```

## 6. Congratulations!
