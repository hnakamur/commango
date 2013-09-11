package executil

import (
	"os/exec"
	"testing"
)

type FormatCommandTestCase struct {
	args []string
	expected string
}

func TestFormatCommand(t *testing.T) {
	cases := []FormatCommandTestCase{
		{[]string{`touch`, `"foo bar"`}, `touch "foo bar"`}, // No need to escape
		{[]string{`touch`, `foo bar`}, `touch foo\\ bar`}, // Escape a space with a backslash
		{[]string{`touch`, `"'foo 'bar"`}, `touch "'foo 'bar"`}, // No need to escape
		{[]string{`touch`, `"foo 'bar"`}, `touch "foo 'bar"`},
		{[]string{`touch`, `'foo "bar'`}, `touch 'foo "bar'`},
		{[]string{`touch`, `'foo ''bar'`}, `touch 'foo ''bar'`},
		{[]string{`touch`, `'foo ' 'bar'`}, `touch 'foo '\\ 'bar'`},
		{[]string{`touch`, `"foo ''bar'"`}, `touch "foo ''bar'"`},
		{[]string{`touch`, `"foo ' 'bar'"`}, `touch "foo ' 'bar'"`},
		{[]string{`touch`, `'foo \'bar'`}, `touch 'foo \'bar'`},
	}
	for _, c := range(cases) {
		cmd := exec.Command(c.args[0], c.args[1:]...)
		actual, err := FormatCommand(cmd)
		if err != nil {
			t.Fatal(err)
		}
		if actual != c.expected {
			t.Errorf("got %s, expected %s.", actual, c.expected)
		}
	}
}
