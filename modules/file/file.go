package file

import (
	"errors"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"

	"github.com/hnakamur/commango/modules"
	"github.com/hnakamur/group"
)

func EnsureLchown(path, owner string) (result modules.Result, err error) {
	extra := make(map[string]interface{})
	extra["op"] = "lchown"
	extra["path"] = path
	extra["owner"] = owner
	result.Extra = extra

	result.RecordStartTime()
	defer func() {
		result.RecordEndTime()

		if err != nil {
			result.Err = err
			result.Failed = true
		}
		result.Log()
		modules.ExitOnError(err)
	}()

	oldUid, oldGid, err := getUidGid(path)
	if err == nil {
		var username, groupname string
		index := strings.IndexAny(owner, ".:")
		if index != -1 {
			username = owner[:index]
			groupname = owner[index+1:]
		} else {
			username = owner
		}

		uid := -1
		gid := -1
		if username != "" {
			extra["user"] = username
			var u *user.User
			u, err = user.Lookup(username)
			if err != nil {
				return
			}

			uid, err = strconv.Atoi(u.Uid)
			if err != nil {
				return
			}
			extra["uid"] = uid
		}

		if groupname != "" {
			extra["group"] = groupname
			var g *group.Group
			g, err = group.Lookup(groupname)
			if err != nil {
				return
			}

			gid, err = strconv.Atoi(g.Gid)
			if err != nil {
				return
			}
			extra["gid"] = gid
		}

		if (uid != -1 && uid != oldUid) ||
			(gid != -1 && gid != oldGid) {

			if uid != -1 && uid != oldUid {
				extra["old_uid"] = oldUid
				var oldUser *user.User
				oldUser, err = user.LookupId(strconv.Itoa(oldUid))
				if err != nil {
					return
				}
				extra["old_user"] = oldUser.Username
			}

			if gid != -1 && gid != oldGid {
				extra["old_gid"] = oldGid
				var oldGroup *group.Group
				oldGroup, err = group.LookupId(strconv.Itoa(oldGid))
				if err != nil {
					return
				}
				extra["old_group"] = oldGroup.Name
			}

			err = os.Lchown(path, uid, gid)
			if err == nil {
				result.Changed = true
			}
		}
	}
	return
}

func getUidGid(path string) (uid, gid int, err error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return -1, -1, err
	}

	stat := fi.Sys().(*syscall.Stat_t)
	if stat == nil {
		return -1, -1, errors.New("not implemented on this platform")
	}

	return int(stat.Uid), int(stat.Gid), nil
}
