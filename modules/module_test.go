package modules

import (
	"testing"
	"time"

	"github.com/hnakamur/commango/jsonutil"
)

func TestRun(t *testing.T) {
	var cm CommandModule
	result, err := Run(cm, "uname", "-a")
	if err != nil {
		t.Fatal(err)
	}

	_, err = jsonutil.Encode(result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEncodeJson(t *testing.T) {
	start := time.Now()
	time.Sleep(123 * time.Millisecond)
	end := time.Now()
	result := Result{
		Failed:    false,
		Changed:   true,
		Err:       nil,
		StartTime: start,
		EndTime:   end,
	}
	_, err := jsonutil.Encode(result)
	if err != nil {
		t.Fatal(err)
	}
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
		actual := FormatDuration(c.d)
		if actual != c.expected {
			t.Fatalf("got %s, expected %s", actual, c.expected)
		}
	}
}
