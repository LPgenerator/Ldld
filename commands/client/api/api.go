package api

import (
	"net/http"

	"github.com/codegangsta/negroni"
	log "github.com/Sirupsen/logrus"

	"github.com/LPgenerator/Ldld/helpers"
	"github.com/LPgenerator/Ldld/commands/server"
	"github.com/LPgenerator/Ldld/commands/client"
)


type LdlCliApi struct{
	path        string
	repo        string
	addr        string
	dist        string
	fs          string
}

var LD struct {
	srv         *server.LdlSrv
	cli         *client.LdlCli
}


func HandleCreate(w http.ResponseWriter, r *http.Request) {
	helpers.DoubleVarView(w, r, "template", "name", LD.cli.Create)
}

func HandleStart(w http.ResponseWriter, r *http.Request) {
	helpers.StandardView(w, r, LD.srv.Start)
}

func HandleStop(w http.ResponseWriter, r *http.Request) {
	helpers.StandardView(w, r, LD.srv.Stop)
}

func HandleInfo(w http.ResponseWriter, r *http.Request) {
	helpers.StandardView(w, r, LD.srv.Info, "LxcInfoToInterface")
}

func HandleList(w http.ResponseWriter, r *http.Request) {
	helpers.JSONReponse(w, LD.srv.List(), "LxcListToInterface")
}

func HandleImages(w http.ResponseWriter, r *http.Request) {
	// todo: to struct list
	helpers.JSONReponse(w, LD.cli.Images())
}

func HandlePull(w http.ResponseWriter, r *http.Request) {
	helpers.StandardView(w, r, LD.cli.Pull)
}

func HandleDestroy(w http.ResponseWriter, r *http.Request) {
	helpers.StandardView(w, r, LD.cli.Destroy)
}

func HandleFreeze(w http.ResponseWriter, r *http.Request) {
	helpers.StandardView(w, r, LD.srv.Freeze)
}

func HandleUnfreeze(w http.ResponseWriter, r *http.Request) {
	helpers.StandardView(w, r, LD.srv.Unfreeze)
}

func HandleExec(w http.ResponseWriter, r *http.Request) {
	helpers.DoubleVarView(w, r, "name", "command", LD.srv.Exec)
}

func HandleLog(w http.ResponseWriter, r *http.Request) {
	helpers.StandardView(w, r, LD.srv.Log, "LxcOutToList")
}

func HandleAutoStart(w http.ResponseWriter, r *http.Request) {
	helpers.DoubleVarView(w, r, "name", "value", LD.cli.Autostart)
}

func HandleForward(w http.ResponseWriter, r *http.Request) {
	helpers.DoubleVarView(w, r, "name", "value", LD.cli.Forward)
}

func HandleCgroup(w http.ResponseWriter, r *http.Request) {
	//helpers.DoubleVarView(w, r, "name", "value", LD.cli.Cgroup)
}

func HandleMemory(w http.ResponseWriter, r *http.Request) {
	helpers.DoubleVarView(w, r, "name", "value", LD.cli.Memory)
}

func HandleSwap(w http.ResponseWriter, r *http.Request) {
	helpers.DoubleVarView(w, r, "name", "value", LD.cli.Swap)
}

func HandleIp(w http.ResponseWriter, r *http.Request) {
	helpers.DoubleVarView(w, r, "name", "value", LD.cli.Ip)
}

func HandleCpu(w http.ResponseWriter, r *http.Request) {
	helpers.DoubleVarView(w, r, "name", "value", LD.cli.Cpu)
}

func HandleProcesses(w http.ResponseWriter, r *http.Request) {
	helpers.DoubleVarView(w, r, "name", "value", LD.cli.Processes)
}

func HandleNetworking(w http.ResponseWriter, r *http.Request) {
	helpers.DoubleVarView(w, r, "name", "value", LD.cli.Networking)
}

func HandleDisk(w http.ResponseWriter, r *http.Request) {
	helpers.DoubleVarView(w, r, "name", "value", LD.cli.Disk)
}

func New(path string, repo string, addr string, dist string, fs string) (*LdlCliApi) {
	strm := &LdlCliApi{
		path: path,
		repo: repo,
		addr: addr,
		dist: dist,
		fs:   fs,
	}
	return strm
}

func (c *LdlCliApi) Run(login string, password string) {
	LD.cli = client.New(c.path, c.repo, c.fs)
	LD.srv = server.New(c.path, c.dist, "", c.fs)

	mux := http.NewServeMux()
	mux.HandleFunc("/", helpers.HandleApiIndex)
	mux.HandleFunc("/create", HandleCreate)
	mux.HandleFunc("/start", HandleStart)
	mux.HandleFunc("/stop", HandleStop)
	mux.HandleFunc("/info", HandleInfo)
	mux.HandleFunc("/list", HandleList)
	mux.HandleFunc("/images", HandleImages)
	mux.HandleFunc("/pull", HandlePull)
	mux.HandleFunc("/destroy", HandleDestroy)
	mux.HandleFunc("/freeze", HandleFreeze)
	mux.HandleFunc("/unfreeze", HandleUnfreeze)
	mux.HandleFunc("/exec", HandleExec)
	mux.HandleFunc("/log", HandleLog)
	mux.HandleFunc("/autostart", HandleAutoStart)
	mux.HandleFunc("/forward", HandleForward)
	mux.HandleFunc("/cgroup", HandleCgroup)
	mux.HandleFunc("/memory", HandleMemory)
	mux.HandleFunc("/swap", HandleSwap)
	mux.HandleFunc("/ip", HandleIp)
	mux.HandleFunc("/cpu", HandleCpu)
	mux.HandleFunc("/processes", HandleProcesses)
	mux.HandleFunc("/network", HandleNetworking)
	mux.HandleFunc("/disk", HandleDisk)

	n := negroni.New()
	n.Use(helpers.LogMiddleware())
	n.Use(helpers.AuthMiddleware(login, password))

	log.Println("Ldl Cli Api listen at", c.addr)
	n.UseHandler(mux)
	if err := http.ListenAndServe(c.addr, n); err != nil {
		log.Errorf("Api exited with error: %s", err)
	}
}
