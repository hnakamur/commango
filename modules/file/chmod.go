package file

import (
	"fmt"
	"os"
	"strings"

	"github.com/hnakamur/commango/stringutil"
	"github.com/hnakamur/commango/task"
)

type Chmod struct {
	Path      string
	Mode      os.FileMode
	Recursive bool
}

func (c *Chmod) Run() (result *task.Result, err error) {
	oldModes, err := c.getModes()
	if err != nil {
		return
	}

	result, err = task.DoRun(func(result *task.Result) (err error) {
		result.Module = "chmod"
		result.Op = "chmod"
		modeStr := fmt.Sprintf("%o", c.Mode)
		result.Extra["path"] = c.Path
		result.Extra["mode"] = modeStr
		result.Extra["recursive"] = c.Recursive
		result.Extra["old_modes"] = oldModes

		if len(oldModes) == 1 {
			oldMode := oldModes[0]

			if modeStr == oldMode {
				result.Skipped = true
				return
			}
		}

		var args []string
		if c.Recursive {
			args = []string{"-R", modeStr, c.Path}
		} else {
			args = []string{modeStr, c.Path}
		}
		err = result.ExecCommand("chmod", args...)
		return
	})
	return
}

func (c *Chmod) getModes() (modes []string, err error) {
	_, err = task.DoRun(func(result *task.Result) (err error) {
		result.Module = "chmod"
		result.Op = "check_modes"
		result.Extra["path"] = c.Path
		result.Extra["recursive"] = c.Recursive

		var args []string
		if c.Recursive {
			args = []string{c.Path, "-printf", "%m\\n"}
		} else {
			args = []string{c.Path, "-printf", "%m\\n", "-quit"}
		}
		err = result.ExecCommand("find", args...)
		result.Changed = false

		lines := strings.Split(strings.TrimRight(result.Stdout, "\n"), "\n")
		modes = stringutil.Uniq(lines)
		return
	})
	return
}
