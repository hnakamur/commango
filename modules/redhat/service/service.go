package service

import (
	"strings"

	"github.com/hnakamur/commango/modules"
)

const (
	STARTED = "started"
	STOPPED = "stopped"
)

func Status(name string) (status string, err error) {
	result, err := modules.CommandNoLog("service", name, "status")
	if result.Rc == 3 {
		status = STOPPED
		result.Err = nil
		err = nil
		result.Failed = false
	} else if result.Rc == 0 {
		status = STARTED
	}
	result.Changed = false
	result.Log()
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

func AutoStartEnabled(name string) (enabled bool, err error) {
	result, err := modules.CommandNoLog("chkconfig", name, "--list")
	enabled = strings.Contains(result.Stdout, "\t2:on\t")
	result.Changed = false
	result.Log()
	return
}

func EnableAutoStart(name string) (result modules.Result, err error) {
	return modules.Command("chkconfig", name, "on")
}

func DisableAutoStart(name string) (result modules.Result, err error) {
	return modules.Command("chkconfig", name, "off")
}
