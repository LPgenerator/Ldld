package helpers

import (
	"os"
	"fmt"
	"strings"
	"io/ioutil"
	"crypto/rand"
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
	hostname, err := os.Hostname(); if err == nil {
		name = fmt.Sprintf("%s.%s", name, hostname)
	}
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


func GenerateMacAddress(mac string) (string, error) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	buf[0] |= 2
	mac_st := strings.Split(mac, ":")
	new_mac := []string{}
	for i, part := range mac_st {
		if part == "xx" {
			part = fmt.Sprintf("%02x", buf[i])
		}
		new_mac = append(new_mac, part)
	}
	return strings.Join(new_mac, ":"), nil
}


func GetDefaultConfig(config string) string {
	lines := strings.Split(config, "\n")
	new_config := ""
	for _, line := range(lines) {
		if strings.HasPrefix(line, "lxc.network.hwaddr") {
			line_split := strings.Split(line, "=")
			curr_mac := strings.Trim(line_split[1], " ")
			new_mac, err := GenerateMacAddress(curr_mac)
			if err == nil {
				line = fmt.Sprintf("lxc.network.hwaddr = %s", new_mac)
			}
		}
		new_config += line + "\n"
	}
	return new_config
}
