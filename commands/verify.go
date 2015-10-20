package commands

import (
	"fmt"
	"encoding/json"

	"github.com/codegangsta/cli"
	log "github.com/Sirupsen/logrus"

	"github.com/LPgenerator/Ldld/common"
)

type VerifyCommand struct {
	configOptions
}

func (c *VerifyCommand) Execute(context *cli.Context) {
	err := c.loadConfig()
	if err != nil {
		log.Fatalln(err)
		return
	}
	data, _ := json.MarshalIndent(c.config, "", "  ")
	fmt.Println(string(data))
}

func init() {
	common.RegisterCommand2("verify", "Verify configuration", &VerifyCommand{})
}
