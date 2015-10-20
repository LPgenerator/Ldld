package helpers

import (
	"os"
	"github.com/codegangsta/cli"
	log "github.com/Sirupsen/logrus"
)

func SetupLogLevelOptions(app *cli.App) {
	newFlags := []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "Debug mode",
			EnvVar: "DEBUG",
		},
		cli.StringFlag{
			Name:  "log-level, l",
			Value: "error",
			Usage: "Log level (options: debug, info, warn, error, fatal, panic)",
		},
	}
	app.Flags = append(app.Flags, newFlags...)
	appBefore := app.Before
	app.Before = func(c *cli.Context) error {
		log.SetOutput(os.Stderr)
		level, err := log.ParseLevel(c.String("log-level"))
		if err != nil {
			log.Fatalf(err.Error())
		}
		log.SetLevel(level)

		if !c.IsSet("log-level") && !c.IsSet("l") && c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		if appBefore != nil {
			return appBefore(c)
		} else {
			return nil
		}
	}
}
