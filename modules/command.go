package modules

import (
	"os/exec"

	"github.com/hnakamur/commango/os/executil"
)

func Command(arg ...string) (result Result, err error) {
	result.RecordStartTime()
	defer result.RecordEndTime()

	cmd := exec.Command(arg[0], arg[1:]...)
	r, err := executil.Run(cmd)

	extra := make(map[string]interface{})
	result.Extra = extra
	extra["cmd"], _ = executil.FormatCommand(cmd)
	if err == nil || executil.IsExitError(err) {
		extra["rc"] = r.Rc
		extra["stdout"] = r.Out.String()
		extra["stderr"] = r.Err.String()
		result.Failed = r.Rc != 0
	} else {
		result.Failed = true
		result.Err = err
	}
	result.Changed = true
	return result, err
}
