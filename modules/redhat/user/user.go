package user

import (
	osuser "os/user"
	"strconv"
	"strings"

	unixgroup "github.com/hnakamur/commango/os/unix/group"
	unixuser "github.com/hnakamur/commango/os/unix/user"
	"github.com/hnakamur/commango/stringutil"
	"github.com/hnakamur/commango/task"
)

type State string

const (
	Present = State("present")
	Absent  = State("absent")
)

const AutoUID = -1

type User struct {
	State       State
	Name        string
	Uid         int
	System      bool     // whether or not system account
	Group       string   // primary group
	Groups      []string // supplementary groups
	Appends     bool     // whether or not appending supplementary groups
	Comment     string
	HomeDir     string
	Shell       string
	RemovesHome bool
}

func (u *User) Run() (result *task.Result, err error) {
	result, err = task.DoRun(func(result *task.Result) (err error) {
		result.Module = "user"
		result.Extra["state"] = string(u.State)
		result.Extra["name"] = u.Name
		if u.State == Absent {
			result.Op = "remove"
			result, err = u.ensureAbsent(result)
		} else {
			result.Extra["uid"] = u.Uid
			result.Extra["system"] = u.System
			result.Extra["group"] = u.Group
			result.Extra["groups"] = u.Groups
			result.Extra["u.Appends"] = u.Appends
			result.Extra["comment"] = u.Comment
			result.Extra["home_dir"] = u.HomeDir
			result.Extra["shell"] = u.Shell
			result, err = u.ensurePresent(result)
		}
		return
	})
	return
}

func (u *User) ensurePresent(result *task.Result) (*task.Result, error) {
	var uidStr string
	if u.Uid != AutoUID {
		uidStr = strconv.Itoa(u.Uid)
	}

	var command string
	args := make([]string, 0)
	var err error
	oldUser, err := unixuser.Lookup(u.Name, nil)
	if err != nil {
		if _, ok := err.(osuser.UnknownUserError); !ok {
			return result, err
		}

		result.Op = "create"
		command = "useradd"
	} else {
		var allGroups []*unixgroup.Group
		allGroups, err = unixgroup.AllGroups()
		if err != nil {
			return result, err
		}

		var oldGroup *unixgroup.Group
		oldGroup, err = unixgroup.LookupId(oldUser.Gid, allGroups)
		if err != nil {
			return result, err
		}

		var oldGroups []string
		oldGroups, err = unixgroup.SupplementaryGroups(u.Name, allGroups)
		if err != nil {
			return result, err
		}

		var groupsWillChange bool
		if u.Appends {
			groupsWillChange = !stringutil.ArrayContainsAll(oldGroups, u.Groups)
		} else {
			groupsWillChange = !stringutil.SetEqual(oldGroups, u.Groups)
		}

		if (uidStr == "" || uidStr == oldUser.Uid) &&
			(u.Group == "" || u.Group == oldUser.Gid || u.Group == oldGroup.Name) &&
			!groupsWillChange &&
			(u.Comment == "" || u.Comment == oldUser.Name) &&
			(u.HomeDir == "" || u.HomeDir == oldUser.HomeDir) &&
			(u.Shell == "" || u.Shell == oldUser.Shell) {
			result.Skipped = true
			return result, err
		}

		if uidStr != "" && uidStr != oldUser.Uid {
			result.Extra["old_uid"], err = strconv.Atoi(oldUser.Uid)
			if err != nil {
				return result, err
			}
		}

		if u.Group != "" && u.Group != oldUser.Gid && u.Group != oldGroup.Name {
			result.Extra["old_gid"], err = strconv.Atoi(oldUser.Gid)
			if err != nil {
				return result, err
			}
		}

		if groupsWillChange {
			result.Extra["old_u.Groups"] = oldGroups
		}

		result.Op = "modify"
		command = "usermod"

		if u.Appends && len(u.Groups) > 0 {
			args = append(args, "-a")
		}
	}
	if uidStr != "" {
		args = append(args, "-u", uidStr)
	}
	if u.Group != "" {
		args = append(args, "-g", u.Group)
	}
	if len(u.Groups) > 0 {
		args = append(args, "-G", strings.Join(u.Groups, ","))
	}
	if u.System {
		args = append(args, "-r")
	}
	if u.Comment != "" {
		args = append(args, "-c", u.Comment)
	}
	if u.HomeDir != "" {
		args = append(args, "-d", u.HomeDir)
	}
	if u.Shell != "" {
		args = append(args, "-s", u.Shell)
	}
	args = append(args, u.Name)

	err = result.ExecCommand(command, args...)
	return result, err
}

func (u *User) ensureAbsent(result *task.Result) (*task.Result, error) {
	var err error
	oldUser, err := unixuser.Lookup(u.Name, nil)
	if err != nil {
		if _, ok := err.(osuser.UnknownUserError); ok {
			err = nil
			result.Skipped = true
		}
		return result, err
	}

	result.Extra["old_uid"], err = strconv.Atoi(oldUser.Gid)
	if err != nil {
		return result, err
	}

	args := make([]string, 0)
	if u.RemovesHome {
		args = append(args, "-r")
	}
	args = append(args, u.Name)

	err = result.ExecCommand("userdel", args...)
	return result, err
}
