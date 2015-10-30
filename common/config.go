package common

import (
	"os"
	"time"
	"bufio"
	"bytes"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)


type Config struct {
	BaseConfig
	ModTime                         time.Time `json:"-"`
	Loaded                          bool      `json:"-"`
}

type BaseConfig struct {
	LdlApiAddress                   string  `toml:"api-address"`
	LdlApiLogin                     string  `toml:"api-login"`
	LdlApiPassword                  string  `toml:"api-password"`

	LdlWebAddress                   string  `toml:"web-address"`
	LdlWebLogin                     string  `toml:"web-login"`
	LdlWebPassword                  string  `toml:"web-password"`

	LdlType                         string  `toml:"type"`

	LdlRepo                         string  `toml:"cli-repo-url"`
	LdlSrvPath                      string  `toml:"srv-path"`
	LdlCliPath                      string  `toml:"cli-data-dir"`
	LdlDist                         string  `toml:"lxc-distro"`
	LdlFS                           string  `toml:"lxc-fs"`
}

func NewConfig() *Config {
	return &Config{
		BaseConfig: BaseConfig{
			LdlApiAddress: "127.0.0.1:9090",
			LdlApiLogin: "ldl",
			LdlApiPassword: "7eNQ4iWLgDw4Q6w",

			LdlWebAddress: "127.0.0.1:9191",
			LdlWebLogin: "admin",
			LdlWebPassword: "7eNQ4iWLgDw4Q6w",

			LdlType: "client",
			LdlDist: "ubuntu",
			LdlRepo: "192.168.0.1",
			LdlSrvPath: "/usr/share/nginx/html",
			LdlCliPath: "/usr/local/var/lib/ldl",
			LdlFS: "overlayfs",
		},
	}
}

func (c *Config) StatConfig(configFile string) error {
	_, err := os.Stat(configFile)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) LoadConfig(configFile string) error {
	info, err := os.Stat(configFile)

	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	if _, err = toml.DecodeFile(configFile, &c.BaseConfig); err != nil {
		return err
	}

	c.ModTime = info.ModTime()
	c.Loaded = true
	return nil
}

func (c *Config) SaveConfig(configFile string) error {
	var newConfig bytes.Buffer
	newBuffer := bufio.NewWriter(&newConfig)

	if err := toml.NewEncoder(newBuffer).Encode(&c.BaseConfig); err != nil {
		log.Fatalf("Error encoding TOML: %s", err)
		return err
	}

	if err := newBuffer.Flush(); err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(configFile), 0700)

	err := ioutil.WriteFile(configFile, newConfig.Bytes(), 0600)
	if err != nil {
		return err
	}

	c.Loaded = true
	return nil
}
