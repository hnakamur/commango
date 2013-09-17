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

func NewResult(module string) *Result {
	return &Result{
		Module: module,
		Extra:  make(map[string]interface{}),
	}
}

func ExecCommand(module string, cmd *exec.Cmd) (result *Result, err error) {
	result = NewResult(module)
	result.RecordStartTime()
	defer result.RecordEndTime()

    result.Command, err = executil.FormatCommand(cmd)
    if err != nil {
        return
    }

    r, err := executil.Run(cmd)
    result.SetExecResult(&r, err)
    return
}

func (r *Result) RecordStartTime() {
	r.StartTime = time.Now()
}

func (r *Result) RecordEndTime() {
	r.EndTime = time.Now()
}

func (r *Result) Delta() time.Duration {
	return r.EndTime.Sub(r.StartTime)
}

func (r *Result) SetExecResult(result *executil.Result, err error) {
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

func (r *Result) Log() {
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

func (r *Result) ToJSON() map[string]interface{} {
	obj := make(map[string]interface{})
	if r.Extra != nil {
		for k, v := range r.Extra {
			obj[k] = v
		}
	}
	obj["module"] = r.Module
	obj["skipped"] = r.Skipped
	obj["failed"] = r.Failed
	obj["changed"] = r.Changed
	if r.Err != nil {
		obj["err"] = r.Err.Error()
	}
	obj["start"] = FormatTime(r.StartTime)
	obj["end"] = FormatTime(r.EndTime)
	obj["delta"] = FormatDuration(r.Delta())
	if r.Command != "" {
		obj["cmd"] = r.Command
		obj["rc"] = r.Rc
		obj["stdout"] = r.Stdout
		obj["stderr"] = r.Stderr
	}
	return obj
}

func (r Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.ToJSON())
}

const TIME_FORMAT = "2006-01-02 15:04:05.00000"

func FormatTime(t time.Time) string {
	return t.Format(TIME_FORMAT)
}

func FormatDuration(d time.Duration) string {
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
