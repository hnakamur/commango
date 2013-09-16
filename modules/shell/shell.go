package shell

import (
    "encoding/json"
    "os/exec"

	log "github.com/cihub/seelog"
    "github.com/hnakamur/commango/modules"
    "github.com/hnakamur/commango/os/executil"
    "github.com/hnakamur/commango/os/osutil"
    "github.com/hnakamur/commango/jsonutil"
)

type Shell struct {
	Shell   string
	Command string
	Chdir   string
	Creates string
    Result  *modules.Result
}

const DEFAULT_SHELL = "/bin/sh"

func (s *Shell) Run() error {
    s.Result = modules.NewResult()
    s.Result.RecordStartTime()
    defer s.Log()
    defer s.Result.RecordEndTime()

    if s.Creates != "" && osutil.Exists(s.Creates) {
        s.Result.Skipped = true
        return nil
    }

    var shell string
    if s.Shell == "" {
        shell = DEFAULT_SHELL
    } else {
        shell = s.Shell
    }

    var command string
    if s.Chdir != "" {
        command = "cd " + executil.QuoteWord(s.Chdir) + "; " + s.Command
    } else {
        command = s.Command
    }
	cmd := exec.Command(shell, "-c", command)
    var err error
    s.Result.Command, err = executil.FormatCommand(cmd)
	if err != nil {
		return err
	}

    r, err := executil.Run(cmd)
    s.Result.SetExecResult(&r, err)
    return err
}

func (s *Shell) MarshalJSON() ([]byte, error) {
    obj := s.Result.ToJSON()
    obj["name"] = "shell"
    obj["command"] = s.Command
    if s.Chdir != "" {
        obj["chdir"] = s.Chdir
    }
    if s.Creates != "" {
        obj["creates"] = s.Creates
    }
	return json.Marshal(obj)
}

func (s *Shell) Log() {
	json, err := jsonutil.Encode(s)
	if err != nil {
		log.Error(err)
	}

	if s.Result.Failed {
		log.Error(json)
	} else if s.Result.Changed {
		log.Info(json)
	} else if !s.Result.Skipped {
		log.Debug(json)
	} else {
		log.Trace(json)
	}
}
