package client

import (
	"os"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"net/http"
	"io/ioutil"

	"github.com/LPgenerator/Ldld/helpers"
	"github.com/LPgenerator/Ldld/helpers/backends"
)

type LdlCli struct{
	path      string
	repo      string
	fs        string
	backend   backends.Fs
}

var CT_TEMPLATE = `
lxc.include = /usr/share/lxc/config/ubuntu.common.conf

lxc.rootfs = /var/lib/lxc/%s/rootfs
lxc.mount = /var/lib/lxc/%s/fstab
lxc.utsname = %s
lxc.arch = amd64
`

var CT_IFCONFIG = `
auto lo
  iface lo inet loopback

iface eth0 inet static
  address %s
  netmask 255.255.0.0
  broadcast 10.1.255.255
  gateway 10.0.3.1
  dns-nameservers 8.8.8.8
  dns-nameservers 8.8.4.4
`

var (
	PORT_FORWARD = `iptables -t nat -A PREROUTING -p tcp --dport %s -j DNAT --to %s:%s`
	RM_PORT_FORWARD = `iptables -t nat -L PREROUTING -n -v --line-numbers| grep %s| awk '{print $1}'| xargs iptables -t nat -D PREROUTING`
	WGET = `wget -c --retry-connrefused -t 0 %s/%s/%s -O %s/%s/%s`
	LinksRegexp = regexp.MustCompile(`">(.*?)/?</a>`)
	DESTROY_CT = `zfs destroy -rR lpg/lxc/%s >&/dev/null; lxc-destroy -f -n %s`
	MIGRATE_CFG = `scp /var/lib/lxc/%s/config %s:/var/lib/lxc/%s/`
	MIGRATE_ZFS = `zfs send lpg/lxc/%s@migrate | ssh %s zfs recv -F lpg/lxc/%s`
	MIGRATE_MP = `ssh %s zfs set mountpoint=/var/lib/lxc/%s/rootfs lpg/lxc/%s`
)


func New(path string, repo string, fs string) (*LdlCli) {
	strm := &LdlCli{
		path:     path,
		repo:     repo,
		fs:       fs,
		backend:  backends.New(fs),
	}
	return strm
}


// ## LXC IMPLEMENTATION ## //
func (c *LdlCli) Create(template string, name string) map[string]string {
	if template == "" || name == "" {
		return c.errorMsg("Template or Name is not set!")
	}

	number := ""
	if strings.Contains(template, ":") {
		data := strings.Split(template, ":")
		template = data[0]
		number = "snap" + data[1]
	}

	err := os.MkdirAll("/var/lib/lxc/" + name, 0755)
	if err != nil {
		return c.errorMsg("Can not create directory")
	}

	default_config, err := ioutil.ReadFile("/etc/lxc/default.conf")
	ct_template := ""
	if err != nil {
		ct_template = CT_TEMPLATE + `
lxc.network.type = veth
lxc.network.flags = up
lxc.network.link = lxcbr0
		`
	} else {
		df_cfg := helpers.GetDefaultConfig(string(default_config))
		ct_template = CT_TEMPLATE + df_cfg
	}

	filename := fmt.Sprintf("/var/lib/lxc/%s/config", name)
	data := fmt.Sprintf(ct_template, name, name, name)
	err = ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		return c.errorMsg("Can not write config")
	}

	cfg_file := fmt.Sprintf("/var/lib/lxc/%s/fstab", name)
	if ioutil.WriteFile(cfg_file, []byte(""), 0644) != nil {
		return c.errorMsg("Can create fstab file")
	}

	res := c.backend.GetSnapshotByTemplate(template, number)
	if res["status"] != "ok" {
		return c.errorMsg("Can not get snapshots")
	}

	res = c.backend.Clone(res["message"], name)
	if res["status"] != "ok" {
		return c.errorMsg("Can not clone fs")
	}

	res = c.backend.Mount(name)
	if res["status"] != "ok" {
		return c.errorMsg("Can not set mount point")
	}

	res = c.backend.Optimize(name)
	if res["status"] != "ok" {
		return c.errorMsg("Can optimze fs")
	}

	post_create := c.backend.AfterCreate(template, name)
	if post_create["status"] != "ok" {
		return post_create
	}

	return map[string]string{"status": "ok", "message": "success"}
}

func (c *LdlCli) Destroy(name string) map[string]string  {
	//todo: c.Stop(name)
	res := helpers.ExecRes("cat /var/lib/lxc/%s/config |grep ipv4|awk '{print $3}'", name)
	if res["status"] == "ok" && res["message"] != "" {
		helpers.ExecRes(RM_PORT_FORWARD, res["message"])
		helpers.ExecRes("service iptables-persistent save")
	}
	return helpers.ExecRes(DESTROY_CT, name, name)
}


// ## REPOSITORY CLIENT IMPLEMENTATION ## //
func (c *LdlCli) Autostart(name string, value string) map[string]string {
	if !helpers.SaveLXCDirective(name, "lxc.start.auto", value) {
		return c.errorMsg("Can not update config")
	}
	return map[string]string{"status": "ok", "message": "success"}
}

func (c *LdlCli) Ip(name string, value string) map[string]string {
	// todo: move default gw to config
	if value == "" {
		return c.errorMsg("VM is not running")
	}
	if value == "fix" {
		res := c.getIP(name)
		value = strings.Trim(res["message"], " ")
		if res["status"] != "ok" {
			return res
		}
	}
	// save to lxc config
	if !helpers.SaveLXCDirective(name, "lxc.network.ipv4", value) {
		return c.errorMsg("Can not update config")
	}
	if !helpers.SaveLXCDirective(name, "lxc.network.ipv4.gateway", "auto") {
		return c.errorMsg("Can not update config")
	}

	// dns
	ct_etc := fmt.Sprintf("/var/lib/lxc/%s/rootfs/etc", name)
	dns_cfg := fmt.Sprintf("%s/resolvconf/resolv.conf.d/original", ct_etc)
	dns_res := helpers.ExecRes("echo 'nameserver 8.8.8.8' > %s", dns_cfg)
	if dns_res["status"] != "ok" {
		return c.errorMsg("Can not write ct dns config")
	}

	// ifconfig
	ct_iface_file := fmt.Sprintf("%s/network/interfaces", ct_etc)
	if ioutil.WriteFile(ct_iface_file, []byte(fmt.Sprintf(CT_IFCONFIG, value)), 0644) != nil {
		return c.errorMsg("Can not write ct ip config")
	}

	return map[string]string{"status": "ok", "message": "success"}
}

func (c *LdlCli) Forward(name string, value string) map[string]string {
	if value == "" {
		return c.errorMsg("Ports not set!")
	}
	res := c.Ip(name, "fix")
	if res["status"] != "ok" {
		return res
	}
	res = c.getIP(name)
	ip := strings.Trim(res["message"], "")
	data := strings.Split(value, ":")
	if len(data) != 2 {
		return c.errorMsg("Error ports format")
	}
	res = helpers.ExecRes(PORT_FORWARD, data[0], ip, data[1])
	helpers.ExecRes("service iptables-persistent save")
	return res
}

func (c *LdlCli) Memory(name string, value string) map[string]string {
	val := strings.Split(value, ":")
	limit := "soft"
	if len(val) > 1 {
		value = val[0]
		limit = val[1]
	}
	helpers.SaveLXCDirective(name, "lxc.cgroup.memory.soft_limit_in_bytes", "0")
	helpers.SaveLXCDirective(name, "lxc.cgroup.memory.limit_in_bytes", "0")
	size := c.convertMbToBytes(value)
	if limit == "soft" {
		return c.doCGroup(name, "memory.soft_limit_in_bytes", size)
	}
	//todo: with swap?
	return c.doCGroup(name, "memory.limit_in_bytes", size)
}

func (c *LdlCli) Swap(name string, value string) map[string]string {
	return c.doCGroup(name, "memory.memsw.limit_in_bytes", c.convertMbToBytes(value))
}

func (c *LdlCli) Cpu(name string, value string) map[string]string {
	return c.doCGroup(name, "cpu.shares", value)
}

func (c *LdlCli) Processes(name string, value string) map[string]string {
	return c.doCGroup(name, "pids.max", value)
}

func (c *LdlCli) Cgroup(name string, group string, value string) map[string]string {
	return c.doCGroup(name, group, value)
}

func (c *LdlCli) Images() map[string]string {
	data := ""

	os.MkdirAll(c.path, 0755)

	local, _ := c.getLocalFiles(c.path)
	data += "Local:\n"
	for _, f := range local {
		data += fmt.Sprintf("\t%s\n", f)
	}

	remote, err := c.getRemoteFiles(c.repo)
	data += "\nRemote:\n"
	if err == nil {
		for _, d := range remote {
			data += fmt.Sprintf("\t%s\n", d)
		}
	} else {
		data += fmt.Sprintf("\terror:%s\n", err.Error())
	}
	return map[string]string{"status": "ok", "message": data}
}

func (c *LdlCli) Pull(dist string) map[string]string {
	data := ""
	remote, err := c.getRemoteFiles(fmt.Sprintf("%s/%s/", c.repo, dist))
	errMake := os.MkdirAll(c.path + "/" + dist, 0755)

	if err == nil && errMake == nil {
		for _, d := range(remote) {
			data += fmt.Sprintf("\t%s\n", d)
			helpers.ExecRes(WGET, c.repo, dist, d, c.path, dist, d)
		}

		return c.importFromPath(dist, remote)
	} else {
		return c.errorMsg(err.Error())
	}
}

func (c *LdlCli) Import(dist string) map[string]string {
	local, _ := c.getLocalFiles(fmt.Sprintf("%s/%s", c.path, dist))
	if len(local) > 0 {
		return c.importFromPath(dist, local)
	}
	return map[string]string{"status": "ok", "message": "success"}
}

/*
func (c *LdlCli) Migrate(name string, ssh string) map[string]string {
	// todo: move with mounted zfs datasets and with CT base images
	res := helpers.ExecRes("ssh %s 'mkdir -p /var/lib/lxc/%s/'", ssh, name)
	if res["status"] != "ok" {
		return res
	}

	//all mounted folders
	//cat /var/lib/lxc/web-3/config|grep lxc.mount.entry|awk '{print $3}'

	//get zfs parent
	//zfs get origin lpg/lxc/web-3|tail -1|awk '{print $3}'|cut -d'/' -f3

	//get snapshots
	//zfs list -t snapshot|grep web@snap

	//get base fs
	//zfs get origin lpg/lxc/web-3|tail -1|awk '{print $3}'|cut -d'@' -f1

	//todo: copy base fs
	//todo: copy snapshots
	//todo: copy images (rsync /usr/local/var/lib/ldl/web/)

	res = helpers.ExecRes(MIGRATE_CFG, name, ssh, name)
	if res["status"] != "ok" {
		return res
	}

	res = helpers.ExecRes(MIGRATE_ZFS, name, ssh, name)
	if res["status"] != "ok" {
		return res
	}

	res = helpers.ExecRes(MIGRATE_MP, ssh, name, name)
	if res["status"] != "ok" {
		return res
	}
	return map[string]string{"status": "ok", "message": "success"}
}
*/

func (c *LdlCli) Migrate(name string, ssh string) map[string]string {
	// todo: rename to MigrateCT
	// Migrate CT with inside Data
	res := helpers.ExecRes("ssh %s 'mkdir -p /var/lib/lxc/%s/'", ssh, name)
	if res["status"] != "ok" {
		return res
	}

	res = helpers.ExecRes(MIGRATE_CFG, name, ssh, name)
	if res["status"] != "ok" {
		return res
	}

	res = helpers.ExecRes("zfs destroy lpg/lxc/%s@migrate", name)
	if res["status"] != "ok" {
		return res
	}

	res = helpers.ExecRes("zfs snapshot lpg/lxc/%s@migrate", name)
	if res["status"] != "ok" {
		return res
	}

	res = helpers.ExecRes(MIGRATE_ZFS, name, ssh, name)
	if res["status"] != "ok" {
		return res
	}

	res = helpers.ExecRes(MIGRATE_MP, ssh, name, name)
	if res["status"] != "ok" {
		return res
	}
	return map[string]string{"status": "ok", "message": "success"}
}

func (c *LdlCli) Mount(name string, src string, dst string) map[string]string {
	// todo: online mount using by: mount -o bind
	if dst != "" && strings.HasPrefix(dst, "/") {
		dst = strings.Replace(dst, "/", "", 1)
	}

	lxc_dir := fmt.Sprintf("/var/lib/lxc/%s", name)
	ovs_delta := fmt.Sprintf("/var/lib/lxc/%s/delta0", name)
	if helpers.FileExists(ovs_delta) {
		// overlayfs
		if os.MkdirAll(fmt.Sprintf("%s/delta0/%s", lxc_dir, dst), 0755) != nil {
			return c.errorMsg("Can not create mount point")
		}
	} else {
		if os.MkdirAll(fmt.Sprintf("%s/rootfs/%s", lxc_dir, dst), 0755) != nil {
			return c.errorMsg("Can not create mount point")
		}
	}

	cfg_file := fmt.Sprintf("%s/fstab", lxc_dir)
	mount_cfg := fmt.Sprintf("%s %s none bind 0 0", src, dst)

	config, _ := ioutil.ReadFile(cfg_file)
	new_config := string(config)
	new_config += fmt.Sprintf("\n%s\n", mount_cfg)

	if ioutil.WriteFile(cfg_file, []byte(new_config), 0644) != nil {
		return c.errorMsg("Can not update fstab")
	}
	return map[string]string{"status": "ok", "message": "success"}
}

func (c *LdlCli) Unmount(name string, src string) map[string]string {
	cfg_file := fmt.Sprintf("/var/lib/lxc/%s/fstab", name)
	config, err := ioutil.ReadFile(cfg_file)
	if err != nil {
		return c.errorMsg("fstab not found")
	}
	new_config := ""
	for _, line := range strings.Split(string(config), "\n") {
		if line != "" {
			line_split := strings.Split(line, " ")
			if src == line_split[0] || strings.HasPrefix(line_split[0], src) {
				continue
			}
			new_config += line + "\n"
		}
	}
	if ioutil.WriteFile(cfg_file, []byte(new_config), 0644) != nil {
		return c.errorMsg("Can not update fstab")
	}
	return map[string]string{"status": "ok", "message": "success"}
}

func (c *LdlCli) Fstab(name string) map[string]string {
	cfg_file := fmt.Sprintf("/var/lib/lxc/%s/fstab", name)
	config, err := ioutil.ReadFile(cfg_file)
	if err != nil {
		return c.errorMsg("fstab not found")
	}
	data := ""
	for _, line := range strings.Split(string(config), "\n") {
		if line != "" {
			line_split := strings.Split(line, " ")
			data += fmt.Sprintf("%s -> /%s", line_split[0], line_split[1])
		}
	}
	return map[string]string{"status": "ok", "message": data}
}


/////
/////
/////


func (c *LdlCli) getIP(name string) map[string]string {
	return helpers.ExecRes("lxc-info -n %s|grep IP|awk '{print $2}'", name)
}

func (c *LdlCli) importFromPath(dist string, path []string) map[string]string {
	for i, _ := range(path) {
		if chk := c.backend.SnapshotIsExists(dist, i); chk["status"] == "ok" && chk["message"] != "" {
			continue
		}

		fmt.Println(fmt.Sprintf("> snap%d", i))
		if res := c.backend.ImportImage(c.path, dist, i); res["status"] != "ok" {
			return res
		}
	}
	return map[string]string{"status": "ok", "message": "success"}
}

func (c *LdlCli) errorMsg(message string) map[string]string {
	return map[string]string{"status": "error", "message": message}
}

func (c *LdlCli) getRemoteFiles(uri string) ([]string, error) {
	response, err := http.Get(uri)
	data := []string{}
	if err == nil {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err == nil {
			dl := LinksRegexp.FindAllStringSubmatch(string(contents), -1)
			for _, d := range dl {
				if d[1] != ".." {
					data = append(data, d[1])
				}
			}
		}
	}
	if err != nil {
		return data, err
	}
	return data, nil
}

func (c *LdlCli) getLocalFiles(path string) ([]string, error) {
	data := []string{}

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		data = append(data, f.Name())
	}
	return data, nil
}

func (c *LdlCli) doCGroup(name string, group string, value string) map[string]string {
	if name == "" || group == "" || value == "" {
		return c.errorMsg("All values must be set!")
	}
	res := helpers.ExecRes("lxc-cgroup -n %s %s %s", name, group, value)
	if res["status"] != "ok" {
		return res
	}
	if !helpers.SaveLXCDirective(name, "lxc.cgroup." + group, value) {
		return c.errorMsg("Can not update config")
	}
	return res
}


func (c *LdlCli) convertMbToBytes(value string) string {
	num, err := strconv.Atoi(value)
	if err == nil {
		val := strconv.Itoa(num * 1024 * 1024)
		return val
	}
	return "0"
}
