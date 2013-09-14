package yum

import (
	"strings"

	"github.com/hnakamur/commango/modules"
)

func Installed(name string) (installed bool, err error) {
	result, err := modules.CommandNoLog("rpm", "-q", name)
	if result.Rc == 0 {
		installed = true
	} else if result.Rc == 1 {
		installed = false
		result.Err = nil
		err = nil
		result.Failed = false
	}
	result.Changed = false
	result.Log()
	modules.ExitOnError(err)
	return
}

func Install(name string) (result modules.Result, err error) {
	result, err = modules.Command("yum", "install", "-d", "2", "-y", name)
	if strings.Contains(result.Stdout, "\nNothing to do\n") {
		result.Changed = false
	}
	return
}
