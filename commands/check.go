package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/LPgenerator/Ldld/common"
)

type CheckCommand struct {
	configOptions
}

func (c *CheckCommand) Execute(context *cli.Context) {
	// todo: check lxc, zfs is installed
	// todo: run lxc-checkconfig
	// todo: check 'lpg/lxc' is exists
	// todo: rename on all code 'lpg/lxc' to 'ldl/lxc'
	// todo: check share dir is exists
	// todo: check images dir is exists
	fmt.Println("All is OK")
}

func init() {
	common.RegisterCommand2(
		"check", "Check installation environment", &CheckCommand{})
}
