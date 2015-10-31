package providers

import (
	"github.com/LPgenerator/Ldld/helpers/providers/lxc"
)

// Solaris: http://habrahabr.ru/post/123221/

// FreeBSD: https://www.opennet.ru/base/sec/jail_example.txt.html
//          http://www.ignix.ru/book/freebsd/setup/jail
//          http://wiki.lissyara.su/wiki/FreeBSD_Руководство_Пользователя:_Глава_15_Jails_(Клетки)

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


	/*
	if provider == "jail" {
		fs = []Provider{jail.Jail{}}
	}
	if provider == "zone" {
		fs = []Provider{zone.Zone{}}
	}
	if provider == "ovz" {
		fs = []Provider{ovz.Ovz{}}
	}
	*/
	return pr[0]
}
