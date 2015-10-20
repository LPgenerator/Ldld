package commands

import (
	"os"
	"path/filepath"

	"github.com/LPgenerator/Ldld/common"
	"github.com/LPgenerator/Ldld/helpers"
)

type configOptions struct {
	config *common.Config

	ConfigFile string `short:"c" long:"config" env:"CONFIG_FILE" description:"Config file"`
}

func getDefaultConfigFile() string {
	if os.Getuid() == 0 {
		return "/etc/ldld/config.toml"
	} else if homeDir := helpers.GetHomeDir(); homeDir != "" {
		return filepath.Join(homeDir, ".ldld", "config.toml")
	} else if currentDir := helpers.GetCurrentWorkingDirectory(); currentDir != "" {
		return filepath.Join(currentDir, "config.toml")
	} else {
		panic("Cannot get default config file location")
	}
}

func (c *configOptions) saveConfig() error {
	return c.config.SaveConfig(c.ConfigFile)
}

func (c *configOptions) loadConfig() error {
	config := common.NewConfig()
	err := config.LoadConfig(c.ConfigFile)
	if err != nil {
		return err
	}
	c.config = config
	return nil
}

func (c *configOptions) touchConfig() error {
	err := c.loadConfig()
	if err != nil {
		return err
	}

	if !c.config.Loaded {
		return c.saveConfig()
	}
	return nil
}

func init() {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		os.Setenv("CONFIG_FILE", getDefaultConfigFile())
	}
}
