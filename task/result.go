package task

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	log "github.com/cihub/seelog"
	"github.com/hnakamur/commango/jsonutil"
	"github.com/hnakamur/commango/os/executil"
)

type Result struct {
	Module    string
	Op        string
	Skipped   bool
	Failed    bool
	Changed   bool
	Err       error
	StartTime time.Time
	EndTime   time.Time
	Command   string
	Rc        int
	Stdout    string
	Stderr    string
	Extra     map[string]interface{}
}

func DoRun(fn func(*Result) error) (result *Result, err error) {
	result = &Result{
		Extra: make(map[string]interface{}),
	}
	result.StartTime = time.Now()
	err = fn(result)
	result.Err = err
	result.EndTime = time.Now()
	result.log()
	return
}

func (r *Result) ExecCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	var err error
	r.Command, err = executil.FormatCommand(cmd)
	if err != nil {
		return err
	}

	runResult, err := executil.Run(cmd)
	r.setExecResult(&runResult, err)
	return err
}

func (r *Result) setExecResult(result *executil.Result, err error) {
	r.Err = err
	if err == nil || executil.IsExitError(err) {
		r.Rc = result.Rc
		r.Stdout = result.Out.String()
		r.Stderr = result.Err.String()
		r.Failed = result.Rc != 0
	} else {
		r.Failed = true
	}
	r.Skipped = false
	r.Changed = true
}

func (r *Result) log() {
	json, err := jsonutil.Encode(r)
	if err != nil {
		log.Error(err)
	}

	if r.Failed {
		log.Error(json)
	} else if r.Changed {
		log.Info(json)
	} else if !r.Skipped {
		log.Debug(json)
	} else {
		log.Trace(json)
	}
}

func (r Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.toJSON())
}

func (r *Result) toJSON() map[string]interface{} {
	obj := make(map[string]interface{})
	if r.Extra != nil {
		for k, v := range r.Extra {
			obj[k] = v
		}
	}
	obj["module"] = r.Module
	obj["op"] = r.Op
	obj["skipped"] = r.Skipped
	obj["failed"] = r.Failed
	obj["changed"] = r.Changed
	if r.Err != nil {
		obj["err"] = r.Err.Error()
	}
	obj["start"] = formatTime(r.StartTime)
	obj["end"] = formatTime(r.EndTime)
	obj["delta"] = formatDuration(r.delta())
	if r.Command != "" {
		obj["command"] = r.Command
		obj["rc"] = r.Rc
		obj["stdout"] = r.Stdout
		obj["stderr"] = r.Stderr
	}
	return obj
}

func (r *Result) delta() time.Duration {
	return r.EndTime.Sub(r.StartTime)
}

const timeFormat = "2006-01-02 15:04:05.00000"

func formatTime(t time.Time) string {
	return t.Format(timeFormat)
}

func formatDuration(d time.Duration) string {
	sign := ""
	if d < 0 {
		d = -d
		sign = "-"
	}

	hour := d / time.Hour
	d = d - hour*time.Hour
	minute := d / time.Minute
	d = d - minute*time.Minute
	second := float64(d) / float64(time.Second)
	return fmt.Sprintf("%s%d:%02d:%08.5f", sign, hour, minute, second)
}
