package modules

import (
	"encoding/json"
	"fmt"
	"time"
)

type Result struct {
	Failed    bool
	Changed   bool
	Err       error
	StartTime time.Time
	EndTime   time.Time
	Extra     map[string]interface{}
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

func (r *Result) ToJSON() map[string]interface{} {
	obj := make(map[string]interface{})
	if r.Extra != nil {
		for k, v := range(r.Extra) {
			obj[k] = v
		}
	}
	obj["failed"] = r.Failed
	obj["changed"] = r.Changed
	if r.Err != nil {
		obj["err"] = r.Err.Error()
	}
	obj["start"] = FormatTime(r.StartTime)
	obj["end"] = FormatTime(r.EndTime)
	obj["delta"] = FormatDuration(r.Delta())
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
