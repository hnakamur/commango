package file

import (
	"strings"

	"github.com/hnakamur/commango/modules"
	"github.com/hnakamur/commango/modules/command"
	"github.com/hnakamur/commango/stringutil"
)

func Chown(path, owner string, recursive bool) (result modules.Result, err error) {
	oldOwner, err := getOwner(path, recursive)
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

	if len(oldOwner) == 1 {
		index := strings.Index(oldOwner[0], ":")
		oldUsername := oldOwner[0][:index]
		oldGroupname := oldOwner[0][index+1:]

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

func getOwner(path string, recursive bool) ([]string, error) {
	var args []string
	if recursive {
		args = []string{"find", path, "-printf", "%u:%g\\n"}
	} else {
		args = []string{"find", path, "-printf", "%u:%g\\n", "-quit"}
	}
	result, err := command.CommandNoLog(args...)
	result.Changed = false
	if err != nil {
		result.Err = err
		result.Failed = true
	}
	result.Log()
	modules.ExitOnError(err)
	owners := strings.Split(strings.TrimRight(result.Stdout, "\n"), "\n")
	return stringutil.Uniq(owners), err
}
