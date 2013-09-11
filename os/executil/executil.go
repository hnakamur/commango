package executil

import (
	"bytes"
	"errors"
	"os/exec"
	"syscall"
)

func CommandLine(cmd *exec.Cmd) string {
	var line bytes.Buffer
	for i, arg := range(cmd.Args) {
		if i > 0 {
			line.WriteByte(' ')
		}
		line.WriteString(arg)
	}
	return line.String()
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
