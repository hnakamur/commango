package file

import (
	"fmt"
	"os"
	"strings"

	"github.com/hnakamur/commango/modules"
	"github.com/hnakamur/commango/modules/command"
	"github.com/hnakamur/commango/stringutil"
)

func Chmod(path string, mode os.FileMode, recursive bool) (result modules.Result, err error) {
	oldModes, err := getModes(path, recursive)
	if err != nil {
		return
	}

	modeStr := fmt.Sprintf("%o", mode)

	result.RecordStartTime()
	defer func() {
		extra := make(map[string]interface{})
		extra["op"] = "chown"
		extra["path"] = path
		extra["mode"] = modeStr
		extra["old_modes"] = oldModes
		result.Extra = extra

		result.RecordEndTime()

		if err != nil {
			result.Err = err
			result.Failed = true
		}
		result.Log()
		modules.ExitOnError(err)
	}()

	if len(oldModes) == 1 {
		oldMode := oldModes[0]

		if mode == oldMode {
			return
		}
	}

	if recursive {
		result, err = command.CommandNoLog("chmod", "-R", modeStr, path)
	} else {
		result, err = command.CommandNoLog("chmod", modeStr, path)
	}
	if err != nil {
		return
	}

	result.Changed = true
	return
}

func getModes(path string, recursive bool) ([]os.FileMode, error) {
	var args []string
	if recursive {
		args = []string{"find", path, "-printf", "%m\\n"}
	} else {
		args = []string{"find", path, "-printf", "%m\\n", "-quit"}
	}
	result, err := command.CommandNoLog(args...)
	result.Changed = false
	if err != nil {
		result.Err = err
		result.Failed = true
	}
	result.Log()
	modules.ExitOnError(err)

	lines := stringutil.Uniq(strings.Split(strings.TrimRight(result.Stdout, "\n"), "\n"))
	modes := make([]os.FileMode, len(lines))
	for i, line := range(lines) {
		modes[i], err = scanMode(line)
		if err != nil {
			break
		}
	}
	return modes, err
}

func scanMode(str string) (mode os.FileMode, err error) {
	n, err := fmt.Sscanf(str, "%o", &mode)
	if err == nil && n != 1 {
		err = fmt.Errorf("scanMode should get one octal number, but got %d numbers from %s",
			n, str)
	}
	return
}