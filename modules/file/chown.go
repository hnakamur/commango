package file

import (
	"strings"

	"github.com/hnakamur/commango/modules"
	"github.com/hnakamur/commango/modules/command"
)

func Chown(path, owner string, recursive bool) (result modules.Result, err error) {
	oldOwner, err := getOwner(path)
	if err != nil {
		return
	}

	result.RecordStartTime()
	defer func() {
		extra := make(map[string]interface{})
		extra["op"] = "chown"
		extra["path"] = path
		extra["owner"] = owner
		extra["old_owner"] = oldOwner
		result.Extra = extra

		result.RecordEndTime()

		if err != nil {
			result.Err = err
			result.Failed = true
		}
		result.Log()
		modules.ExitOnError(err)
	}()

	index := strings.Index(oldOwner, ":")
	oldUsername := oldOwner[:index]
	oldGroupname := oldOwner[index+1:]

	var username, groupname string
	index = strings.IndexAny(owner, ".:")
	if index != -1 {
		username = owner[:index]
		groupname = owner[index+1:]
	} else {
		username = owner
	}

	if (username == "" || username == oldUsername) &&
		(groupname == "" || groupname == oldGroupname) {
		return
	}

	if recursive {
		result, err = command.CommandNoLog("chown", "-R", owner, path)
	} else {
		result, err = command.CommandNoLog("chown", owner, path)
	}
	if err != nil {
		return
	}

	result.Changed = true
	return
}

func getOwner(path string) (string, error) {
	result, err := command.CommandNoLog("find", path, "-printf", "%u:%g", "-quit")
	result.Changed = false
	if err != nil {
		result.Err = err
		result.Failed = true
	}
	result.Log()
	modules.ExitOnError(err)
	return result.Stdout, err
}
