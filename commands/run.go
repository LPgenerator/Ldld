package commands

import (
	"os"

	"github.com/codegangsta/cli"
	log "github.com/Sirupsen/logrus"
	"github.com/LPgenerator/Ldld/common"
	"github.com/LPgenerator/Ldld/helpers/service"
	service "github.com/ayufan/golang-kardianos-service"

	cli_api "github.com/LPgenerator/Ldld/commands/client/api"
	srv_api "github.com/LPgenerator/Ldld/commands/server/api"
)


type RunCommand struct {
	configOptions

	ListenAddr       string `short:"l" long:"listen" description:"Listen address:port"`
	ServiceName      string `short:"n" long:"service" description:"Use different names for different services"`
	WorkingDirectory string `short:"d" long:"working-directory" description:"Specify custom working directory"`
	Syslog           bool   `long:"syslog" description:"Log to syslog"`
}


func (mr *RunCommand) Run() {
	if mr.config.LdlType == "client" {
		api := cli_api.New(mr.config.LdlCliPath, mr.config.LdlRepo, mr.config.LdlApiAddress, mr.config.LdlDist)
		api.Run(mr.config.LdlApiLogin, mr.config.LdlApiPassword)
	} else {
		api := srv_api.New(mr.config.LdlSrvPath, mr.config.LdlRepo, mr.config.LdlApiAddress, mr.config.LdlDist)
		api.Run(mr.config.LdlApiLogin, mr.config.LdlApiPassword)
	}
}

func (mr *RunCommand) Start(s service.Service) error {
	if len(mr.WorkingDirectory) > 0 {
		err := os.Chdir(mr.WorkingDirectory)
		if err != nil {
			return err
		}
	}

	err := mr.loadConfig()
	if err != nil {
		panic(err)
	}

	go mr.Run()

	return nil
}

func (mr *RunCommand) Stop(s service.Service) error {
	log.Println("Ldld: requested service stop")
	// mr.saveConfig()
	return nil
}

func (c *RunCommand) Execute(context *cli.Context) {
	svcConfig := &service.Config{
		Name:        c.ServiceName,
		DisplayName: c.ServiceName,
		Description: defaultDescription,
		Arguments:   []string{"run"},
	}

	service, err := service_helpers.New(c, svcConfig)
	if err != nil {
		log.Fatalln(err)
	}

	if c.Syslog {
		logger, err := service.SystemLogger(nil)
		if err == nil {
			log.AddHook(&ServiceLogHook{logger})
		} else {
			log.Errorln(err)
		}
	}

	err = service.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	common.RegisterCommand2("run", "Run Ldl http API", &RunCommand{
		ServiceName: defaultServiceName,
	})
}
