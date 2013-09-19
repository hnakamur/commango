package shell

import (
    "github.com/hnakamur/commango/os/executil"
    "github.com/hnakamur/commango/os/osutil"
    "github.com/hnakamur/commango/task"
)

type Shell struct {
	Shell   string
	Command string
	Chdir   string
	Creates string
}

const DEFAULT_SHELL = "/bin/sh"

func (s *Shell) Run() (result *task.Result, err error) {
    result = task.NewResult("shell")
    result.RecordStartTime()
    defer func() {
        result.RecordEndTime()
        result.Log()
    }()

    result.Extra["command"] = s.Command
    if s.Chdir != "" {
        result.Extra["chdir"] = s.Chdir
    }
    if s.Creates != "" {
        result.Extra["creates"] = s.Creates
    }

    var shell string
    if s.Shell == "" {
        shell = DEFAULT_SHELL
    } else {
        shell = s.Shell
    }
    result.Extra["shell"] = shell

    if s.Creates != "" && osutil.Exists(s.Creates) {
        result.Skipped = true
        return
    }

    var command string
    if s.Chdir != "" {
        command = "cd " + executil.QuoteWord(s.Chdir) + "; " + s.Command
    } else {
        command = s.Command
    }
	err = result.ExecCommand(shell, "-c", command)
    return
}
