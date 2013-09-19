package yum

import (
	"strings"

	"github.com/hnakamur/commango/task"
)

type State string

const (
	Installed = State("installed")
	Removed   = State("removed")
)

type Yum struct {
	State State
	Name  string
}

func (y *Yum) Run() (result *task.Result, err error) {
    installed, err := y.isInstalled()
	if err != nil {
		return
	}

	result = task.NewResult("yum")
	result.RecordStartTime()
	result.Extra["name"] = y.Name

	if y.State == Installed {
		if installed {
			result.Skipped = true
		} else {
			err = result.ExecCommand("yum", "install", "-d", "2", "-y", y.Name)
			if strings.Contains(result.Stdout, "\nNothing to do\n") {
				result.Changed = false
			}
		}
	} else {
		if installed {
			err = result.ExecCommand("yum", "-C", "remove", "-y", y.Name)
		} else {
			result.Skipped = true
		}
	}
	result.RecordEndTime()
	result.Log()
	return
}

func (y *Yum) isInstalled() (installed bool, err error) {
	result := task.NewResult("yum.installed")
	result.RecordStartTime()
	result.Extra["name"] = y.Name

	err = result.ExecCommand("rpm", "-q", y.Name)
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
	return
}
