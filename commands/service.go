package commands

import (
	"os"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
	log "github.com/Sirupsen/logrus"
	service "github.com/ayufan/golang-kardianos-service"

	"github.com/LPgenerator/Ldld/common"
	"github.com/LPgenerator/Ldld/helpers"
	"github.com/LPgenerator/Ldld/helpers/service"
)

const (
	defaultServiceName = "ldld"
	defaultDisplayName = "Ldld is a simple and open source tool for shipping and running distributed containers."
	defaultDescription = "Ldld is a simple and open source tool for shipping and running distributed containers."
)

type ServiceLogHook struct {
	service.Logger
}

func (s *ServiceLogHook) Levels() []log.Level {
	return []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
		log.InfoLevel,
	}
}

func (s *ServiceLogHook) Fire(e *log.Entry) error {
	switch e.Level {
	case log.PanicLevel, log.FatalLevel, log.ErrorLevel:
		s.Error(e.Message)
	case log.WarnLevel:
		s.Warning(e.Message)
	case log.InfoLevel:
		s.Info(e.Message)
	}
	return nil
}

type NullService struct {
}

func (n *NullService) Start(s service.Service) error {
	return nil
}

func (n *NullService) Stop(s service.Service) error {
	return nil
}

func runServiceInstall(s service.Service, c *cli.Context) error {
	if user := c.String("user"); user == "" && os.Getuid() == 0 {
		log.Fatal("Please specify user that will run ldld service")
	}

	if configFile := c.String("config"); configFile != "" {
		// try to load existing config
		config := common.NewConfig()
		err := config.LoadConfig(configFile)
		if err != nil {
			return err
		}

		// save config for the first time
		if !config.Loaded {
			err = config.SaveConfig(configFile)
			if err != nil {
				return err
			}
		}
	}
	return service.Control(s, "install")
}

func RunServiceControl(c *cli.Context) {
	// detect whether we want to install as user service or system service
	isUserService := os.Getuid() != 0

	// when installing service as system wide service don't specify username for service
	serviceUserName := c.String("user")
	if !isUserService {
		serviceUserName = ""
	}

	if isUserService && runtime.GOOS == "linux" {
		log.Fatal("Please run the commands as root")
	}

	svcConfig := &service.Config{
		Name:        c.String("service"),
		DisplayName: c.String("service"),
		Description: defaultDescription,
		Arguments:   []string{"run"},
		UserName:    serviceUserName,
	}

	switch runtime.GOOS {
	case "darwin":
		svcConfig.Option = service.KeyValue{
			"KeepAlive":     true,
			"RunAtLoad":     true,
			"SessionCreate": true,
			"UserService":   isUserService,
		}
	}

	if wd := c.String("working-directory"); wd != "" {
		svcConfig.Arguments = append(
			svcConfig.Arguments, "--working-directory", wd)
	}

	if config := c.String("config"); config != "" {
		svcConfig.Arguments = append(svcConfig.Arguments, "--config", config)
	}

	if sn := c.String("service"); sn != "" {
		svcConfig.Arguments = append(svcConfig.Arguments, "--service", sn)
	}

	// svcConfig.Arguments = append(svcConfig.Arguments, "--syslog")

	s, err := service_helpers.New(&NullService{}, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	switch c.Command.Name {
	case "service-install":
		err = runServiceInstall(s, c)
	default:
		err = service.Control(s, strings.Replace(c.Command.Name, "service-", "", -1))
	}

	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "service, n",
			Value: defaultServiceName,
			Usage: "Specify service name to use",
		},
	}

	installFlags := flags
	installFlags = append(installFlags, cli.StringFlag{
		Name:  "working-directory, d",
		Value: helpers.GetCurrentWorkingDirectory(),
		Usage: "Specify custom root directory where all data are stored",
	})
	installFlags = append(installFlags, cli.StringFlag{
		Name:  "config, c",
		Value: getDefaultConfigFile(),
		Usage: "Specify custom config file",
	})

	if runtime.GOOS == "windows" {
		installFlags = append(installFlags, cli.StringFlag{
			Name:  "user, u",
			Value: helpers.GetCurrentUserName(),
			Usage: "Specify user-name to secure the LB",
		})
		installFlags = append(installFlags, cli.StringFlag{
			Name:  "password, p",
			Value: "",
			Usage: "Specify user password to install service (required)",
		})
	} else if os.Getuid() == 0 {
		installFlags = append(installFlags, cli.StringFlag{
			Name:  "user, u",
			Value: "",
			Usage: "Specify user-name to secure the LB",
		})
	}

	common.RegisterCommand(cli.Command{
		Name:   "service-install",
		Usage:  "Install service",
		Action: RunServiceControl,
		Flags:  installFlags,
	})
	common.RegisterCommand(cli.Command{
		Name:   "service-uninstall",
		Usage:  "Uninstall service",
		Action: RunServiceControl,
		Flags:  flags,
	})
	common.RegisterCommand(cli.Command{
		Name:   "service-start",
		Usage:  "Start service",
		Action: RunServiceControl,
		Flags:  flags,
	})
	common.RegisterCommand(cli.Command{
		Name:   "service-stop",
		Usage:  "Stop service",
		Action: RunServiceControl,
		Flags:  flags,
	})
	common.RegisterCommand(cli.Command{
		Name:   "service-restart",
		Usage:  "Restart service",
		Action: RunServiceControl,
		Flags:  flags,
	})
}
