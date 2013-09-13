package modules

import (
	"os/exec"

	"github.com/hnakamur/commango/os/executil"
)

type CommandModule struct {
	dummy string
}

type CommandResult struct {
	*ModuleResult
	Cmd    string `json:"cmd"`
	Rc     int    `json:"rc"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func (c CommandModule) Main(arg ...string) (interface{}, error) {
	cmd := exec.Command(arg[0], arg[1:]...)
	r, err := executil.Run(cmd)

	var mr ModuleResult
	cr := CommandResult{ModuleResult: &mr}
	cr.Cmd, _ = executil.FormatCommand(cmd)
	if err == nil || executil.IsExitError(err) {
		cr.Rc = r.Rc
		cr.Stdout = r.Out.String()
		cr.Stderr = r.Err.String()
		failed := r.Rc != 0
		cr.Failed = failed
	} else {
		cr.Failed = true
		cr.Err = err
	}
	cr.Changed = true
	return cr, err
}
