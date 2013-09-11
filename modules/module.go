package modules

import (
	"fmt"
	"time"
)

type ResultJson map[string]interface{}

type Module interface {
	Main(arg ...string) (ResultJson, error)
}

type Runner struct{}

func (r *Runner) Run(m Module, arg ...string) (ResultJson, error) {
	start := time.Now()

	rj, err := m.Main(arg...)

	end := time.Now()
	delta := end.Sub(start)

	if rj == nil {
		rj = ResultJson{}
	}
	rj["start"] = formatTime(start)
	rj["end"] = formatTime(end)
	rj["delta"] = formatDuration(delta)
	return rj, err
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
	d = d - hour * time.Hour
	minute := d / time.Minute
	d = d - minute * time.Minute
	second := float64(d) / float64(time.Second)
	return fmt.Sprintf("%s%d:%02d:%08.5f", sign, hour, minute, second)
}
