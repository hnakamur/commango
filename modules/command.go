package modules

import (
	"os/exec"

	"github.com/hnakamur/commango/os/executil"
)

func CommandNoLog(arg ...string) (result Result, err error) {
	result.RecordStartTime()
	defer result.RecordEndTime()

	cmd := exec.Command(arg[0], arg[1:]...)
	result.Command, err = executil.FormatCommand(cmd)
	if err != nil {
		return
	}

	r, err := executil.Run(cmd)
	result.SetExecResult(&r, err)
	return
}

func Command(arg ...string) (result Result, err error) {
	result, err = CommandNoLog(arg...)
	result.Log()
	ExitOnError(err)
	return
}
