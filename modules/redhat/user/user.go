package user

import (
	osuser "os/user"
	"strconv"
	"strings"

	"github.com/hnakamur/commango/modules"
	"github.com/hnakamur/commango/modules/command"
	unixgroup "github.com/hnakamur/commango/os/unix/group"
	unixuser "github.com/hnakamur/commango/os/unix/user"
	"github.com/hnakamur/commango/stringutil"
)

func EnsureExists(name string, uid int, system bool, group string,
	groups []string, appends bool, comment string, home string,
	shell string) (result modules.Result, err error) {
	result.RecordStartTime()

	extra := map[string]interface{}{
		"op":      "user",
		"name":    name,
		"uid":     uid,
		"system":  system,
		"group":   group,
		"groups":  groups,
		"appends": appends,
		"comment": comment,
		"home":    home,
		"shell":   shell,
	}

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

	var uidStr string
	if uid != -1 {
		uidStr = strconv.Itoa(uid)
	}

	args := make([]string, 0)
	oldUser, err := unixuser.Lookup(name, nil)
	if err != nil {
		if _, ok := err.(osuser.UnknownUserError); !ok {
			return
		}

		args = append(args, "useradd")
	} else {
		extra["old_uid"], err = strconv.Atoi(oldUser.Uid)
		if err != nil {
			return
		}

		extra["old_gid"], err = strconv.Atoi(oldUser.Gid)
		if err != nil {
			return
		}

		var allGroups []*unixgroup.Group
		allGroups, err = unixgroup.AllGroups()
		if err != nil {
			return
		}

		var oldGroup *unixgroup.Group
		oldGroup, err = unixgroup.LookupId(oldUser.Gid, allGroups)
		if err != nil {
			return
		}

		var oldGroups []string
		oldGroups, err = unixgroup.SupplementaryGroups(name, allGroups)
		if err != nil {
			return
		}
		extra["old_groups"] = oldGroups

		var groupsWillChange bool
		if appends {
			groupsWillChange = !stringutil.ArrayContainsAll(oldGroups, groups)
		} else {
			groupsWillChange = !stringutil.SetEqual(oldGroups, groups)
		}

		if (uidStr == "" || uidStr == oldUser.Uid) &&
			(group == "" || group == oldUser.Gid || group == oldGroup.Name) &&
			!groupsWillChange &&
			(comment == "" || comment == oldUser.Name) &&
			(home == "" || home == oldUser.HomeDir) &&
			(shell == "" || shell == oldUser.Shell) {
			result.Skipped = true
			return
		}

		args = append(args, "usermod")

		if appends && len(groups) > 0 {
			args = append(args, "-a")
		}
	}
	if uidStr != "" {
		args = append(args, "-u", uidStr)
	}
	if group != "" {
		args = append(args, "-g", group)
	}
	if len(groups) > 0 {
		args = append(args, "-G", strings.Join(groups, ","))
	}
	if system {
		args = append(args, "-r")
	}
	if comment != "" {
		args = append(args, "-c", comment)
	}
	if home != "" {
		args = append(args, "-d", home)
	}
	if shell != "" {
		args = append(args, "-s", shell)
	}
	args = append(args, name)

	result, err = command.CommandNoLog(args...)
	result.Changed = true
	return
}

func EnsureRemoved(name string, removesHome bool) (result modules.Result, err error) {
	result.RecordStartTime()

	extra := make(map[string]interface{})
	extra["op"] = "userdel"
	extra["name"] = name
	extra["removes_home"] = removesHome

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

	oldUser, err := unixuser.Lookup(name, nil)
	if err != nil {
		if _, ok := err.(osuser.UnknownUserError); ok {
			err = nil
			result.Skipped = true
		}
		return
	}

	extra["old_uid"], err = strconv.Atoi(oldUser.Gid)
	if err != nil {
		return
	}

	args := []string{"userdel"}
	if removesHome {
		args = append(args, "-r")
	}
	args = append(args, name)

	result, err = command.CommandNoLog(args...)
	result.Changed = true
	return
}
