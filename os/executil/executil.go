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

func Run(c *exec.Cmd, okExitStatuses []int) (exitStatus int, err error) {
	if err = c.Start(); err != nil {
		return
	}
	return Wait(c, okExitStatuses)
}

func Wait(cmd *exec.Cmd, okExitStatuses []int) (exitStatus int, err error) {
	err = cmd.Wait()
	exitStatus = getExitStatus(err)
	if err != nil && isExitStatusOk(exitStatus, okExitStatuses) {
		err = nil
	}
	return
}

func getExitStatus(waitResult error) int {
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

func isExitStatusOk(exitStatus int, okExitStatuses []int) bool {
	if okExitStatuses == nil {
		return exitStatus == 0
	}

	for _, s := range(okExitStatuses) {
		if s == exitStatus {
			return true
		}
	}
	return false
}
