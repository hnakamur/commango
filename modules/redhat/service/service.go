package service

import (
	"github.com/hnakamur/commango/modules"
)

func Status(name string) (result modules.Result, err error) {
	result, err = modules.Command("service", name, "status")
	if result.Rc == 3 {
		result.Err = nil
		err = nil
		result.Failed = false
	}
	result.Changed = false
	return
}

func Start(name string) (result modules.Result, err error) {
	return modules.Command("service", name, "start")
}

func Stop(name string) (result modules.Result, err error) {
	return modules.Command("service", name, "stop")
}

func Restart(name string) (result modules.Result, err error) {
	return modules.Command("service", name, "restart")
}

func Reload(name string) (result modules.Result, err error) {
	return modules.Command("service", name, "reload")
}

func AutoStartEnabled(name string) (result modules.Result, err error) {
	result, err = modules.Command("sh", "-c", "chkconfig "+name+" --list | grep -q 2:on")
	if result.Rc == 1 {
		result.Err = nil
		err = nil
		result.Failed = false
	}
	result.Changed = false
	return
}

func EnableAutoStart(name string) (result modules.Result, err error) {
	return modules.Command("chkconfig", name, "on")
}

func DisableAutoStart(name string) (result modules.Result, err error) {
	return modules.Command("chkconfig", name, "off")
}
