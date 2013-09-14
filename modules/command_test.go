package modules

import (
	"testing"

	"github.com/hnakamur/commango/jsonutil"
)

func TestCommand(t *testing.T) {
	result, err := Command("uname", "-a")
	if err != nil {
		t.Fatal(err)
	}

	_, err = jsonutil.Encode(result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommandError(t *testing.T) {
	_, err := Command("sh", "-c", "exit 1")
	if err == nil {
		t.Fatal("command should failed, but not")
	}
}
