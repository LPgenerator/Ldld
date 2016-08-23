package commands


import (
	"fmt"
	"reflect"

	"github.com/codegangsta/cli"
	"gitlab.com/ayufan/golang-cli-helpers"

	"github.com/LPgenerator/Ldld/common"
	"github.com/LPgenerator/Ldld/commands/server"
	"github.com/LPgenerator/Ldld/commands/client"
)


type RunShell struct {
	configOptions
}

type fn func(name string)

func (r *RunShell) regCmd(fn interface{}, c *cli.Context, cn string, name string, args... interface{}) {
	if c.Command.Name != cn { return }
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	call := reflect.ValueOf(fn).MethodByName(name).Call(inputs)
	result := call[0].Interface().(map[string]string)
	if len(call) > 0 && len(result) > 0 {
		if result["status"] != "ok" {
			fmt.Println("ERROR:")
		}
		if result["message"] != "" {
			fmt.Println(result["message"])
		}
	}
}

func (r *RunShell) RunLdlControlSrv(c *cli.Context) {
	key := c.Args().Get(0)
	val := c.Args().Get(1)

	r.loadConfig()

	if key == "" && c.Command.Name != "list" && c.Command.Name != "share" {
		fmt.Println("CT name is not defined!")
		return
	}
	if val == "" && c.Command.Name == "clone" {
		fmt.Println("Destenation is not defined!")
		return
	}
	srv := server.New(r.config.LdlSrvPath, r.config.LdlDist, r.config.LdlCliPath, r.config.LdlFS)

	r.regCmd(srv, c, "create", "Create", key)
	r.regCmd(srv, c, "start", "Start", key)
	r.regCmd(srv, c, "stop", "Stop", key)
	r.regCmd(srv, c, "attach", "Attach", key)
	r.regCmd(srv, c, "list", "List")
	r.regCmd(srv, c, "freeze", "Freeze", key)
	r.regCmd(srv, c, "unfreeze", "Unfreeze", key)
	r.regCmd(srv, c, "info", "Info", key)
	r.regCmd(srv, c, "commit", "Commit", key)
	r.regCmd(srv, c, "log", "Log", key)
	r.regCmd(srv, c, "push", "Push", key)
	r.regCmd(srv, c, "clone", "Clone", key, val)
	r.regCmd(srv, c, "destroy", "Destroy", key)
	r.regCmd(srv, c, "share", "Share")
	r.regCmd(srv, c, "exec", "Exec", key, val)
	r.regCmd(srv, c, "export", "Export", key, val)
}

func (r *RunShell) RunLdlControlCli(c *cli.Context) {
	key := c.Args().Get(0)
	val := c.Args().Get(1)

	r.loadConfig()

	if key == "" && c.Command.Name != "list" && c.Command.Name != "images" {
		fmt.Println("CT name is not defined!")
		return
	}
	if val == "" && c.Command.Name == "clone" {
		fmt.Println("Destenation is not defined!")
		return
	}
	if (val == "" || c.Args().Get(2) == "") && c.Command.Name == "mount" {
		fmt.Println("Src or Dst is not defined!")
		return
	}
	if (val == "" || c.Args().Get(2) == "") && c.Command.Name == "cgroup" {
		fmt.Println("Key or Rule is not defined!")
		return
	}

	cli := client.New(r.config.LdlCliPath, r.config.LdlRepo, r.config.LdlFS)
	srv := server.New(r.config.LdlSrvPath, r.config.LdlDist, "", r.config.LdlFS)

	r.regCmd(cli, c, "create", "Create", val, key)
	r.regCmd(srv, c, "start", "Start", key)
	r.regCmd(srv, c, "stop", "Stop", key)
	r.regCmd(srv, c, "info", "Info", key)
	r.regCmd(srv, c, "list", "List")
	r.regCmd(cli, c, "images", "Images")
	r.regCmd(cli, c, "pull", "Pull", key)
	r.regCmd(cli, c, "import", "Import", key)
	r.regCmd(cli, c, "migrate", "Migrate", key, val)
	r.regCmd(cli, c, "destroy", "Destroy", key)
	r.regCmd(srv, c, "attach", "Attach", key)
	r.regCmd(srv, c, "freeze", "Freeze", key)
	r.regCmd(srv, c, "unfreeze", "Unfreeze", key)
	r.regCmd(srv, c, "share", "Share")
	r.regCmd(srv, c, "exec", "Exec", key, val)
	r.regCmd(srv, c, "log", "Log", key)

	r.regCmd(cli, c, "ip", "Ip", key, val)
	r.regCmd(cli, c, "memory", "Memory", key, val)
	r.regCmd(cli, c, "swap", "Swap", key, val)
	r.regCmd(cli, c, "cpu", "Cpu", key, val)
	r.regCmd(cli, c, "cgroup", "Cgroup", key, val, c.Args().Get(2))
	r.regCmd(cli, c, "autostart", "Autostart", key, val)
	r.regCmd(cli, c, "forward", "Forward", key, val)
	r.regCmd(cli, c, "processes", "Processes", key, val)

	r.regCmd(cli, c, "mount", "Mount", key, val, c.Args().Get(2))
	r.regCmd(cli, c, "unmount", "Unmount", key, val)
	r.regCmd(cli, c, "fstab", "Fstab", key)
}

func init() {
	srv_commands := map[string]string{
		"create": "Creates a container.",
		"commit": "Snapshot an existing container.",
		"push": "Push container diff to repository.",
		"export": "Export to remote client without repo sharing",
		"share": "Share local repository.",

		"start": "Run container.",
		"stop": "Stop a container.",
		"list": " List the containers existing on the system.",
		"clone": "Clone a new container from an existing one.",
		"freeze": "Freeze all the container's processes.",
		"unfreeze": "Thaw all the container's processes.",
		"info": "Query information about a container.",
		"log": "List an existing snapshots for container.",
		"attach": "Enter into a running container.",
		"exec": "Execute the command inside CT.",
		"destroy": "Destroy a container.",
	}
	cli_commands := map[string]string{
		"images": "List the images existing on the system.",
		"pull": "Pull container updates from repository.",
		"import": "Import when used export on master.",
		"migrate": "Migrate CT to a new host.",
		"create": "Creates a container.",
		"start": "Run container.",
		"stop": "Stop a container.",
		"info": "Query information about a container.",
		"list": " List the containers existing on the system.",
		"destroy": "Destroy a container.",
		"attach": "Enter into a running container.",
		"exec": "Execute the command inside CT.",
		"freeze": "Freeze all the container's processes.",
		"unfreeze": "Thaw all the container's processes.",
		"log": "List an existing snapshots for container.",

		"processes": "Processes limits.",
		"memory": "Memory limits. Set a maximum RAM (on MB).",
		"swap": "Swap limits. Set a maximum swap (on MB).",
		"cpu": "CPU limits based on CPU shares.",
		"cgroup": "",
		"mount": "",
		"unmount": "",
		"fstab": "",
		"autostart": "Autostart after a reboot (0 or 1).",
		"forward": "Port forwarding (ip will be fixed).",
		"ip": "Static IP. If value = 'fix', current IP will be fixed.",
	}

	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "name, n",
			Value: "",
		}, cli.StringFlag{
			Name:  "to, t",
			Value: "",
		},
	}
	data := &RunShell{}
	available_commands := srv_commands
	action := data.RunLdlControlSrv

	config := common.NewConfig()
	err := config.LoadConfig(getDefaultConfigFile())

	if err == nil && config.LdlType == "client" {
		available_commands = cli_commands
		action = data.RunLdlControlCli
	}

	for k, v := range(available_commands) {
		common.RegisterCommand(cli.Command{
			Name:   k,
			Usage:  v,
			Action: action,
			Flags: append(flags, clihelpers.GetFlagsFromStruct(data)...),
		})
	}
}
