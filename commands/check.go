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
	// todo: check lxc, zfs/btrfs/overlayds is installed
	// todo: run lxc-checkconfig
	// todo: check zfs:lpg/lxc || btrfs:/var/lib/lxc is exists
	// todo: check share dir is exists
	// todo: check images dir is exists
	// todo: ... and etc
	fmt.Println("Not yet!")
}

func init() {
	common.RegisterCommand2(
		"check", "Check installation environment", &CheckCommand{})
}
