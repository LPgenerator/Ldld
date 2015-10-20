package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	_ "github.com/LPgenerator/Ldld/commands"
	"github.com/LPgenerator/Ldld/common"
	"github.com/LPgenerator/Ldld/helpers"
)

var NAME = "ldld"
var VERSION = "dev"
var REVISION = "HEAD"

func init() {
	common.NAME = NAME
	common.VERSION = VERSION
	common.REVISION = REVISION
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "ldld"
	app.Version = fmt.Sprintf("%s (%s)", common.VERSION, common.REVISION)
	app.Author = "GoTLiuM InSPiRiT"
	app.Email = "gotlium@gmail.com"
	helpers.SetupLogLevelOptions(app)
	app.Commands = common.GetCommands()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

