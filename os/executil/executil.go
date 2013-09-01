package executil

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"syscall"
)

type CommandRunner struct {
	Command *exec.Cmd
	CapturesStdout bool
	CapturesStderr bool
	stdoutBuffer *bytes.Buffer
	stderrBuffer *bytes.Buffer
	OkExitStatuses []int
	exitStatus int
}

func (r *CommandRunner) Run() error {
	cmd := r.Command
	stdoutLogger := NewLogger(Info)
	if r.CapturesStdout {
		var outBuf bytes.Buffer
		r.stdoutBuffer = &outBuf
		cmd.Stdout = io.MultiWriter(r.stdoutBuffer, stdoutLogger)
	} else {
		cmd.Stdout = stdoutLogger
	}
	stderrLogger := NewLogger(Err)
	if r.CapturesStderr {
		var errBuf bytes.Buffer
		r.stderrBuffer = &errBuf
		cmd.Stderr = io.MultiWriter(r.stderrBuffer, stderrLogger)
	} else {
		cmd.Stderr = stderrLogger
	}

	line := r.CommandLine()
	stdoutLogger.Logf("run command\tcommand:%s", line)
	err := cmd.Run()
	if err != nil {
		if e2, ok := err.(*exec.ExitError); ok {
			if s, ok := e2.Sys().(syscall.WaitStatus); ok {
				r.exitStatus = s.ExitStatus()
			} else {
				panic(errors.New("Unimplemented for system where exec.ExitError.Sys() is not syscall.WaitStatus."))
			}
		}
	} else {
		r.exitStatus = 0
	}
	if r.IsExitStatusOk() {
		err = nil
		stdoutLogger.Logf("done\tstatus:%d", r.exitStatus)
	} else {
		stderrLogger.Logf("failed\tstatus:%d", r.exitStatus)
	}
	return err
}

func (r *CommandRunner) StdoutOutput() string {
	if r.CapturesStdout {
		return r.stdoutBuffer.String()
	} else {
		return ""
	}
}

func (r *CommandRunner) StderrOutput() string {
	if r.CapturesStderr {
		return r.stderrBuffer.String()
	} else {
		return ""
	}
}

func (r *CommandRunner) CommandLine() string {
	cmd := r.Command
	var line bytes.Buffer
	for i, arg := range(cmd.Args) {
		if i > 0 {
			line.WriteByte(' ')
		}
		line.WriteString(arg)
	}
	return line.String()
}

func (r *CommandRunner) IsExitStatusOk() bool {
	if r.OkExitStatuses == nil {
		return r.exitStatus == 0
	}

	for _, s := range(r.OkExitStatuses) {
		if s == r.exitStatus {
			return true
		}
	}
	return false
}

func (r *CommandRunner) ExitStatus() int {
	return r.exitStatus
}
