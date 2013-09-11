package executil

import (
	"bytes"
	"errors"
	"os/exec"
	"syscall"
)

type Result struct {
	Out bytes.Buffer // the stdout output
	Err bytes.Buffer // the stderr output
	Rc  int          // the exit status
}

func Run(cmd *exec.Cmd) (result Result, err error) {
	cmd.Stdout = &result.Out
	cmd.Stderr = &result.Err
	err = cmd.Run()
	result.Rc = GetExitStatus(err)
	return
}

func CommandLine(cmd *exec.Cmd) string {
	var line bytes.Buffer
	for i, arg := range cmd.Args {
		if i > 0 {
			line.WriteByte(' ')
		}
		line.WriteString(arg)
	}
	return line.String()
}

func IsExitError(err error) bool {
	_, ok := err.(*exec.ExitError)
	return ok
}

func GetExitStatus(waitResult error) int {
	if waitResult != nil {
		if err, ok := waitResult.(*exec.ExitError); ok {
			if s, ok := err.Sys().(syscall.WaitStatus); ok {
				return s.ExitStatus()
			} else {
				panic(errors.New("Unimplemented for system where exec.ExitError.Sys() is not syscall.WaitStatus."))
			}
		}
	}
	return 0
}
