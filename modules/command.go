package modules

import (
	"os/exec"

	"github.com/hnakamur/commango/os/executil"
)

type CommandModule struct {
	dummy string
}

func (c *CommandModule) Main(arg ...string) (ResultJson, error) {
	cmd := exec.Command(arg[0], arg[1:]...)
	r, err := executil.Run(cmd)

	rj := ResultJson{}
	rj["cmd"] = executil.CommandLine(cmd)
	if err == nil || executil.IsExitError(err) {
		rj["rc"] = r.Rc
		rj["stdout"] = r.Out.String()
		rj["stderr"] = r.Err.String()
		failed := r.Rc != 0
		rj["failed"] = failed
	} else {
		rj["failed"] = true
		rj["err"] = err.Error()
	}
	rj["changed"] = true
	return rj, err
}
