package file

import (
	"strings"

	"github.com/hnakamur/commango/task"
	"github.com/hnakamur/commango/stringutil"
)

type Chown struct {
	Path      string
	Owner     string
	Group     string
	Recursive bool
}

func (c *Chown) Run() (result *task.Result, err error) {
	oldOwners, err := c.getOwners()
	if err != nil {
		return
	}

	result = task.NewResult("chown")
	result.RecordStartTime()

	result.Extra["path"] = c.Path
	result.Extra["owner"] = c.Owner
	result.Extra["group"] = c.Group
	result.Extra["recursive"] = c.Recursive

	defer func() {
		result.RecordEndTime()
		result.Log()
	}()

	if len(oldOwners) == 1 {
		index := strings.Index(oldOwners[0], ":")
		oldOwner := oldOwners[0][:index]
		oldGroup := oldOwners[0][index+1:]

		if (c.Owner == "" || c.Owner == oldOwner) &&
			(c.Group == "" || c.Group == oldGroup) {
			result.Skipped = true
			return
		}
	}

	var owner string
	if c.Group != "" {
		owner = c.Owner + ":" + c.Group
	} else {
		owner = c.Owner
	}
	var args []string
	if c.Recursive {
		args = []string{"-R", owner, c.Path}
	} else {
		args = []string{owner, c.Path}
	}
	err = result.ExecCommand("chown", args...)
	return
}

func (c *Chown) getOwners() (owners []string, err error) {
	result := task.NewResult("chown.get_owners")
	result.RecordStartTime()

	result.Extra["path"] = c.Path
	result.Extra["recursive"] = c.Recursive

	var args []string
	if c.Recursive {
		args = []string{c.Path, "-printf", "%u:%g\\n", "-quit"}
	} else {
		args = []string{c.Path, "-printf", "%u:%g\\n"}
	}
	err = result.ExecCommand("find", args...)

	lines := strings.Split(strings.TrimRight(result.Stdout, "\n"), "\n")
	owners = stringutil.Uniq(lines)

	result.Changed = false
	result.RecordEndTime()
	result.Log()
	return
}
