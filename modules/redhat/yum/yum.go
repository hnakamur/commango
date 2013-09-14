package yum

import (
	"strings"

	"github.com/hnakamur/commango/modules"
	"github.com/hnakamur/commango/modules/command"
)

func Installed(name string) (installed bool, err error) {
	result, err := command.CommandNoLog("rpm", "-q", name)
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
	result, err = command.Command("yum", "install", "-d", "2", "-y", name)
	if strings.Contains(result.Stdout, "\nNothing to do\n") {
		result.Changed = false
	}
	return
}

func EnsureInstalled(name string) (result modules.Result, err error) {
	installed, err := Installed(name)
	if installed || err != nil {
		return
	}

	return Install(name)
}
