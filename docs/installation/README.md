## Manual installation and configuration

### Install

Simply download one of the binaries for your system:

```bash
wget -O /usr/local/bin/ldld https://github.com/LPgenerator/Ldld/releases/download/v1.0/ldld-linux-386
wget -O /usr/local/bin/ldld https://github.com/LPgenerator/Ldld/releases/download/v1.0/ldld-linux-amd64
wget -O /usr/local/bin/ldld https://github.com/LPgenerator/Ldld/releases/download/v1.0/ldld-linux-arm
```

Give it permissions to execute:

```bash
chmod +x /usr/local/bin/ldld
```

Create a ldld user:

```
useradd --create-home ldld --shell /bin/bash
```

Install and run as service:
```bash
sudo ldld install --user=ldld --working-directory=/home/ldld
sudo ldld start
```

### Update

Stop the service (you need elevated command prompt as before):

```bash
sudo ldld stop
```

Download the binary to replace LB's executable:

```bash
wget -O /usr/local/bin/ldld https://github.com/LPgenerator/Ldld/releases/download/v1.0/ldld-linux-386
wget -O /usr/local/bin/ldld https://github.com/LPgenerator/Ldld/releases/download/v1.0/ldld-linux-amd64
wget -O /usr/local/bin/ldld https://github.com/LPgenerator/Ldld/releases/download/v1.0/ldld-linux-arm
```

Give it permissions to execute:

```bash
chmod +x /usr/local/bin/ldld
```

Start the service:

```bash
ldld start
```

### Config

Config file must be found after installation at the `~/.ldld/config.toml` or in `/etc/ldld/config.toml` 
