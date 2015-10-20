package api

import (
	"net/http"

	"github.com/codegangsta/negroni"
	log "github.com/Sirupsen/logrus"

	"github.com/LPgenerator/Ldld/helpers"
	"github.com/LPgenerator/Ldld/commands/server"
	"github.com/LPgenerator/Ldld/commands/client"
)

type LdlCliSrv struct{
	path        string
	repo        string
	addr        string
	dist        string
}

var LD struct {
	srv         *server.LdlSrv
	cli         *client.LdlCli
}


func HandleCreate(w http.ResponseWriter, r *http.Request) {
	helpers.StandardView(w, r, LD.srv.Create)
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

func HandleClone(w http.ResponseWriter, r *http.Request) {
	helpers.DoubleVarView(w, r, "from", "to", LD.srv.Clone)
}

func HandlePush(w http.ResponseWriter, r *http.Request) {
	helpers.StandardView(w, r, LD.srv.Push)
}
func HandleCommit(w http.ResponseWriter, r *http.Request) {
	helpers.StandardView(w, r, LD.srv.Commit)
}

func New(path string, repo string, addr string, dist string) (*LdlCliSrv) {
	strm := &LdlCliSrv{
		path: path,
		repo: repo,
		addr: addr,
		dist: dist,
	}
	return strm
}

func (c *LdlCliSrv) Run(login string, password string) {
	LD.cli = client.New(c.path, c.repo)
	LD.srv = server.New(c.path, c.dist, "")

	mux := http.NewServeMux()
	mux.HandleFunc("/", helpers.HandleApiIndex)
	mux.HandleFunc("/create", HandleCreate)
	mux.HandleFunc("/start", HandleStart)
	mux.HandleFunc("/stop", HandleStop)
	mux.HandleFunc("/info", HandleInfo)
	mux.HandleFunc("/list", HandleList)
	mux.HandleFunc("/clone", HandleClone)
	mux.HandleFunc("/push", HandlePush)
	mux.HandleFunc("/destroy", HandleDestroy)
	mux.HandleFunc("/freeze", HandleFreeze)
	mux.HandleFunc("/unfreeze", HandleUnfreeze)
	mux.HandleFunc("/exec", HandleExec)
	mux.HandleFunc("/log", HandleLog)
	mux.HandleFunc("/commit", HandleCommit)

	n := negroni.New()
	n.Use(helpers.LogMiddleware())
	n.Use(helpers.AuthMiddleware(login, password))

	log.Println("Ldl Srv Api listen at", c.addr)
	n.UseHandler(mux)
	if err := http.ListenAndServe(c.addr, n); err != nil {
		log.Errorf("Api exited with error: %s", err)
	}
}
