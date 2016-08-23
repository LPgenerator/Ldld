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
gvm install go1.7
gvm use go1.7
gvm pkgset create ldld
gvm pkgset use ldld
```

## 2. Download Ldld sources

```bash
git clone https://github.com/LPgenerator/Ldld.git
mkdir -p ~/.gvm/pkgsets/go1.7/ldld/src/github.com/LPgenerator/
ln -sf `pwd` ~/.gvm/pkgsets/go1.7/ldld/src/github.com/LPgenerator/Ldld
cd ~/.gvm/pkgsets/go1.7/ldld/src/github.com/LPgenerator/Ldld
```

## 3. Install Ldld dependencies

This will download and restore all dependencies required to build `Ldld`:

```bash
make deps
```

## 4. Run Ldld

Normally you would use `Ldld`, in order to compile and run Go source use go toolchain:

```bash
go run main.go run
```

You can run `Ldld` in debug-mode:

```bash
go run main.go --debug run
```

## 5. Compile and install Ldld binary

```bash
go build
go install
```

## 6. Congratulations!
