package group

import (
	"strconv"

	unixgroup "github.com/hnakamur/commango/os/unix/group"
	"github.com/hnakamur/commango/modules"
	"github.com/hnakamur/commango/modules/command"
)

func EnsureExists(name string, gid int, system bool) (result modules.Result, err error) {
	result.RecordStartTime()

	extra := make(map[string]interface{})
	extra["op"] = "group"
	extra["name"] = name
	extra["gid"] = gid
	extra["system"] = system

	defer func() {
		result.Extra = extra

		result.RecordEndTime()

		if err != nil {
			result.Err = err
			result.Failed = true
		}
		result.Log()
		modules.ExitOnError(err)
	}()

	var gidStr string
	if gid != -1 {
		gidStr = strconv.Itoa(gid)
	}

	args := make([]string, 0)
	oldGroup, err := unixgroup.Lookup(name, nil)
	if err != nil {
		if _, ok := err.(unixgroup.UnknownGroupError); !ok {
			return
		}

		args = append(args, "groupadd")
	} else {
		extra["old_gid"], err = strconv.Atoi(oldGroup.Gid)
		if err != nil {
			return
		}
		if gidStr == "" || gidStr == oldGroup.Gid {
			result.Skipped = true
			return
		}

		args = append(args, "groupmod")
	}
	if gidStr != "" {
		args = append(args, "-g", gidStr)
	}
	if system {
		args = append(args, "-r")
	}
	args = append(args, name)

	result, err = command.CommandNoLog(args...)
	result.Changed = true
	return
}

func EnsureRemoved(name string) (result modules.Result, err error) {
	result.RecordStartTime()

	extra := make(map[string]interface{})
	extra["op"] = "groupdel"
	extra["name"] = name

	defer func() {
		result.Extra = extra

		result.RecordEndTime()

		if err != nil {
			result.Err = err
			result.Failed = true
		}
		result.Log()
		modules.ExitOnError(err)
	}()

	oldGroup, err := unixgroup.Lookup(name, nil)
	if err != nil {
		if _, ok := err.(unixgroup.UnknownGroupError); ok {
			err = nil
			result.Skipped = true
		}
		return
	}

	extra["old_gid"], err = strconv.Atoi(oldGroup.Gid)
	if err != nil {
		return
	}

	result, err = command.CommandNoLog("groupdel", name)
	result.Changed = true
	return
}
