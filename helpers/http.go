package helpers

import (
	"io"
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/mailgun/oxy/utils"
	log "github.com/Sirupsen/logrus"
)

type handler func(w http.ResponseWriter, r *http.Request)
type fn func(name string) map[string]string
type fn2 func(v1 string, v2 string) map[string]string

func LogRequests(h handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	}
}

type LMiddleware struct {}
type AMiddleware struct {
	Username string
	Password string
}

func LogMiddleware() *LMiddleware {
	return &LMiddleware{}
}

func (l *LMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
}


func AuthMiddleware(username string, password string) *AMiddleware {
	return &AMiddleware{
		Username: username,
		Password: password,
	}
}

func (a *AMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	auth, err := utils.ParseAuthHeader(r.Header.Get("Authorization"))
	if err != nil || a.Username != auth.Username || a.Password != auth.Password {
		w.Header().Set("WWW-Authenticate", `Basic realm="Auth required"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	next(w, r)
}

func ErrorMSG(w http.ResponseWriter, message string) {
	http.Error(
		w, fmt.Sprintf(`{"status": "error", "message": "%s"}`, message),
		http.StatusInternalServerError)
}

func jsonReponse(w http.ResponseWriter, response interface{}) {
	data, err := json.MarshalIndent(response, "", "  ")
	if err == nil {
		io.WriteString(w, string(data))
		return
	} else {
		ErrorMSG(w, "Cold not encode JSON message")
	}
}

func JSONReponse(w http.ResponseWriter, response map[string]string, formaters ...interface{}) {
	if len(formaters) > 0 && response["status"] == "ok" {
		var format interface{}

		switch formaters[0] {
		case "LxcInfoToInterface":
			format = LxcInfoToInterface(response["message"])
		case "LxcListToInterface":
			format = LxcListToInterface(response["message"])
		default:
			format = LxcOutToList(response["message"])
		}
		jsonReponse(w, map[string]interface{}{
			"status": response["status"],
			"message": format,
		})
		return
	}
	jsonReponse(w, response)
}

func StandardView(w http.ResponseWriter, r *http.Request, f fn, formaters ...interface{}) {
	r.ParseForm()
	name := r.Form.Get("name")
	if r.Method == "POST" && name != ""{
		result := f(name)
		//if result["status"] == "ok" && result["message"] == "" {
		//	result["message"] = "success"
		//}
		if result["status"] == "ok" {
			JSONReponse(w, result, formaters...)
		} else {
			ErrorMSG(w, result["message"])
		}
		return
	}
	ErrorMSG(w, "error")
}

func DoubleVarView(w http.ResponseWriter, r *http.Request, v1 string, v2 string, f fn2, formaters ...interface{}) {
	r.ParseForm()
	v1 = r.Form.Get(v1)
	v2 = r.Form.Get(v2)
	if r.Method == "POST" && v1 != "" && v2 != "" {
		result := f(v1, v2)
		if result["status"] == "ok" {
			JSONReponse(w, result, formaters...)
		} else {
			ErrorMSG(w, result["message"])
		}
		return
	}
	ErrorMSG(w, "error")
}

func HandleApiIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `{"status": "hello"}`)
}
