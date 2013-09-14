package yum

import (
	"strings"

	"github.com/hnakamur/commango/modules"
)

func Installed(name string) (result modules.Result, err error) {
	result, err = modules.Command("rpm", "-q", name)
	if result.Rc == 1 {
        result.Err = nil
        err = nil
    }
	result.Changed = false
	return
}

func Install(name string) (result modules.Result, err error) {
	result, err = modules.Command("yum", "install", "-d", "2", "-y", name)
    if strings.Contains(result.Stdout, "\nNothing to do\n") {
        result.Changed = false
    }
	return
}
