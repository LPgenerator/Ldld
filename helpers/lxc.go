package helpers

import (
	"os"
	"fmt"
	"strings"
	"io/ioutil"
)


func SaveLXCDirective(name string, group string, value string) bool {
	config_filename := fmt.Sprintf("/var/lib/lxc/%s/config", name)
	config, err := ioutil.ReadFile(config_filename)
	if err != nil { return false }
	new_config := ""
	cfg_found := false
	for _, line := range strings.Split(string(config), "\n") {
		if strings.HasPrefix(line, group + " =") {
			if value != "0" {
				new_config += fmt.Sprintf("%s = %s\n", group, value)
				cfg_found = true
			}
		} else {
			if line != "" {
				new_config += line + "\n"
			}
		}
	}
	if cfg_found == false && value != "0" {
		new_config += fmt.Sprintf("%s = %s\n", group, value)
	}
	if ioutil.WriteFile(config_filename, []byte(new_config), 0644) == nil {
		return true
	}
	return false
}


func SaveHostInfo(name string, ct_etc string) map[string]string {
	if !FileExists(ct_etc) && os.MkdirAll(ct_etc, 0755) != nil {
		return map[string]string{"status": "error", "message": "can not create etc"}
	}

	if res := ExecRes("echo %s > %s/hostname", name, ct_etc); res["status"] != "ok" {
		return map[string]string{"status": "error", "message": "can not save hostname"}
	}

	if res := ExecRes("echo '127.0.0.1\t%s' >> %s/hosts", name, ct_etc); res["status"] != "ok" {
		return map[string]string{"status": "error", "message": "can not save hosts"}
	}

	return map[string]string{"status": "ok", "message": "success"}
}
