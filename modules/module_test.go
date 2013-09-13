package modules

import (
	"testing"
	"time"

	"github.com/hnakamur/commango/jsonutil"
)

func TestRun(t *testing.T) {
	var cm CommandModule
	var r Runner
	rj, err := r.Run(cm, "uname", "-a")
	if err != nil {
		t.Fatal(err)
	}

	json, err := jsonutil.Encode(rj)
	t.Fatalf("json=%s.", json)

	//_, err = jsonutil.Encode(rj)
	//if err != nil {
	//    t.Fatal(err)
	//}
}

func TestEncodeJson(t *testing.T) {
	start := time.Now()
	time.Sleep(123 * time.Millisecond)
	end := time.Now()
	delta := end.Sub(start)
	mr := ModuleResult{
		Failed:  false,
		Changed: true,
		Err:     nil,
		Start:   Time(start),
		End:     Time(end),
		Delta:   Duration(delta),
	}
	json, err := jsonutil.Encode(mr)
	if err != nil {
		t.Fatal(err)
	}
	t.Fatalf("json=%s.", json)
}

type FormatDurationTestCase struct {
	d        time.Duration
	expected string
}

func TestFormatDuration(t *testing.T) {
	cases := []FormatDurationTestCase{
		{time.Second + 12*time.Millisecond + 340*time.Microsecond, "0:00:01.01234"},
		{-(time.Second + 12*time.Millisecond + 340*time.Microsecond), "-0:00:01.01234"},
		{987*time.Hour + 59*time.Minute + 12*time.Second + 345670*time.Microsecond, "987:59:12.34567"},
		{987*time.Hour + 59*time.Minute + 12*time.Second + 345675*time.Microsecond, "987:59:12.34567"},
		{987*time.Hour + 59*time.Minute + 12*time.Second + 345676*time.Microsecond, "987:59:12.34568"},
	}
	for _, c := range cases {
		actual := formatDuration(c.d)
		if actual != c.expected {
			t.Fatalf("got %s, expected %s", actual, c.expected)
		}
	}
}
