package providers

import (
	"github.com/LPgenerator/Ldld/helpers/providers/lxc"
)


type Provider interface {
    Create(name string, fs string, dist string) map[string]string
    Start(name string) map[string]string
    Stop(name string) map[string]string
    Attach(name string)
    List() map[string]string
    Freeze(name string) map[string]string
    Unfreeze(name string) map[string]string
    Info(name string) map[string]string
    Clone(src string, dst string) map[string]string
    Destroy(name string) map[string]string
    Exec(name string, cmd string) map[string]string
}


func New(provider string) (Provider) {
	pr := []Provider{lxc.Lxc{}}

	if provider == "ovz" {
		pr = []Provider{ovz.Ovz{}}
	}

	return pr[0]
}
