package group

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Group struct {
	Gid     string // group id
	Name    string
	Members []string // usernames who are members of this group as
	// a supplementary group, not a primary group.
}

// UnknownGroupIdError is returned by LookupId when
// a group cannot be found.
type UnknownGroupIdError int

func (e UnknownGroupIdError) Error() string {
	return "group: unknown groupid " + strconv.Itoa(int(e))
}

// UnknownGroupError is returned by Lookup when
// a group cannot be found.
type UnknownGroupError string

func (e UnknownGroupError) Error() string {
	return "group: unknown group " + string(e)
}

func (g *Group) String() string {
	return "Group[name=" + g.Name + ",gid=" + g.Gid +
		",members=\"" + strings.Join(g.Members, ",") + "\"]"
}

const (
	NAME_COL    = 0
	GID_COL     = 2
	MEMBERS_COL = 3
)

func AllGroups() ([]*Group, error) {
	content, err := ioutil.ReadFile("/etc/group")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	groups := make([]*Group, 0)
	for _, line := range lines {
		if strings.HasPrefix(line, "#") ||
			len(strings.TrimSpace(line)) == 0 {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) <= MEMBERS_COL {
			return nil, fmt.Errorf("Invalid group line: %s", line)
		}
		groups = append(groups, &Group{
			Name: fields[NAME_COL],
			Gid: fields[GID_COL],
			Members: strings.Split(fields[MEMBERS_COL], ","),
		})
	}
	return groups, nil
}

func Lookup(groupname string, groups []*Group) (*Group, error) {
	if groups == nil {
		var err error
		groups, err = AllGroups()
		if err != nil {
			return nil, err
		}
	}
	for _, group := range groups {
		if group.Name == groupname {
			return group, nil
		}
	}
	return nil, UnknownGroupError(groupname)
}

func LookupId(gid string, groups []*Group) (*Group, error) {
	gidInt, err := strconv.Atoi(gid)
	if err != nil {
		return nil, err
	}
	if groups == nil {
		var err error
		groups, err = AllGroups()
		if err != nil {
			return nil, err
		}
	}
	for _, group := range groups {
		if group.Gid == gid {
			return group, nil
		}
	}
	return nil, UnknownGroupIdError(gidInt)
}

func SupplementaryGroups(username string, groups []*Group) ([]string, error) {
	if groups == nil {
		var err error
		groups, err = AllGroups()
		if err != nil {
			return nil, err
		}
	}
	supplementaryGroups := make([]string, 0)
	for _, group := range groups {
		for _, member := range group.Members {
			if member == username {
				supplementaryGroups = append(supplementaryGroups, group.Name)
			}
		}
	}
	return supplementaryGroups, nil
}
