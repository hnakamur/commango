package group

import (
	"strconv"

	unixgroup "github.com/hnakamur/commango/os/unix/group"
	"github.com/hnakamur/commango/task"
)

type State string

const (
	Present = State("present")
	Absent  = State("absent")
)

const AUTO_GID = -1

type Group struct {
	State   State
	Name    string
	Gid     int
	System  bool     // whether or not system account
}

func (g *Group) Run() (result *task.Result, err error) {
	result = task.NewResult("group")
	result.RecordStartTime()

	result.Extra["state"] = string(g.State)
	result.Extra["name"] = g.Name
	if g.State == Absent {
		result, err = g.ensureAbsent(result)
	} else {
		result.Extra["gid"] = g.Gid
		result.Extra["system"] = g.System
		result, err = g.ensurePresent(result)
	}

    result.Err = err
    result.RecordEndTime()
    result.Log()
    return
}

func (g *Group) ensurePresent(result *task.Result) (*task.Result, error) {
	var gidStr string
	if g.Gid != AUTO_GID {
		gidStr = strconv.Itoa(g.Gid)
	}

    var command string
    var err error
	oldGroup, err := unixgroup.Lookup(g.Name, nil)
	if err != nil {
		if _, ok := err.(unixgroup.UnknownGroupError); !ok {
			return result, err
		}

		command = "groupadd"
	} else {
		if gidStr == "" || gidStr == oldGroup.Gid {
			result.Skipped = true
			return result, err
		}

		result.Extra["old_gid"], err = strconv.Atoi(oldGroup.Gid)
		if err != nil {
			return result, err
		}

		command = "groupmod"
	}

	args := make([]string, 0)
	if gidStr != "" {
		args = append(args, "-g", gidStr)
	}
	if g.System {
		args = append(args, "-r")
	}
	args = append(args, g.Name)

    err = result.ExecCommand(command, args...)
	return result, err
}

func (g *Group) ensureAbsent(result *task.Result) (*task.Result, error) {
    var err error
	oldGroup, err := unixgroup.Lookup(g.Name, nil)
	if err != nil {
		if _, ok := err.(unixgroup.UnknownGroupError); ok {
			err = nil
			result.Skipped = true
		}
		return result, err
	}

	result.Extra["old_gid"], err = strconv.Atoi(oldGroup.Gid)
	if err != nil {
		return result, err
	}

    err = result.ExecCommand("groupdel", g.Name)
	return result, err
}
