package modules

import (
	"errors"
	"fmt"
	"time"
)

type ResultJson map[string]interface{}

type ModuleResult struct {
	Failed  bool     `json:"failed"`
	Changed bool     `json:"changed"`
	Err     error    `json:"err",omitempty`
	Start   Time     `json:"start"`
	End     Time     `json:"end"`
	Delta   Duration `json:"delta"`
}

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(t).Format(TIME_FORMAT) + "\""), nil
}

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte("\"" + formatDuration(time.Duration(d)) + "\""), nil
}

type Module interface {
	Main(arg ...string) (interface{}, error)
}

type Runner struct{}

func (r *Runner) Run(m Module, arg ...string) (interface{}, error) {
	start := time.Now()

	result, err := m.Main(arg...)
	mr, ok := result.(ModuleResult)
	if !ok {
		return nil, errors.New("ModeResult expected")
	}

	end := time.Now()
	delta := end.Sub(start)

	mr.Start = Time(start)
	mr.End = Time(end)
	mr.Delta = Duration(delta)
	return mr, err
}

const TIME_FORMAT = "2006-01-02 15:04:05.00000"

func formatTime(t time.Time) string {
	return t.Format(TIME_FORMAT)
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
