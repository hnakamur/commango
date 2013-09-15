package user

import (
	"fmt"
	"io/ioutil"
	osuser "os/user"
	"strconv"
	"strings"
	"syscall"
)

type User struct {
	Uid      string // user id
	Gid      string // primary group id
	Username string
	Name     string
	HomeDir  string
	Shell    string
}

func (u *User) String() string {
	return "User[username=" + u.Username + ",uid=" + u.Uid + ",gid=" + u.Gid +
		",name=" + u.Name + ",homedir=" + u.HomeDir + ",shell=" + u.Shell + "]"
}

const (
	USERNAME_COL = 0
	UID_COL      = 2
	GID_COL      = 3
	NAME_COL     = 4
	HOME_DIR_COL = 5
	SHELL_COL    = 6
)

func AllUsers() ([]*User, error) {
	content, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	users := make([]*User, 0)
	for _, line := range lines {
		if strings.HasPrefix(line, "#") ||
			len(strings.TrimSpace(line)) == 0 {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) <= SHELL_COL {
			return nil, fmt.Errorf("Invalid user line: %s", line)
		}
		users = append(users, &User{
			Uid:      fields[UID_COL],
			Gid:      fields[GID_COL],
			Username: fields[USERNAME_COL],
			HomeDir:  fields[HOME_DIR_COL],
			Shell:    fields[SHELL_COL],
		})
	}
	return users, nil
}

func Current(users []*User) (*User, error) {
	return LookupId(strconv.Itoa(syscall.Getuid()), users)
}

func Lookup(username string, users []*User) (*User, error) {
	if users == nil {
		var err error
		users, err = AllUsers()
		if err != nil {
			return nil, err
		}
	}
	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, osuser.UnknownUserError(username)
}

func LookupId(uid string, users []*User) (*User, error) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return nil, err
	}
	if users == nil {
		var err error
		users, err = AllUsers()
		if err != nil {
			return nil, err
		}
	}
	for _, user := range users {
		if user.Uid == uid {
			return user, nil
		}
	}
	return nil, osuser.UnknownUserIdError(uidInt)
}
