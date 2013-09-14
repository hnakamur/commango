package modules

import (
	"testing"

	"github.com/hnakamur/commango/jsonutil"
)

func TestCommand(t *testing.T) {
	var cm CommandModule
	result, err := cm.Main("uname", "-a")
	if err != nil {
		t.Fatal(err)
	}

	_, err = jsonutil.Encode(result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommandError(t *testing.T) {
	var cm CommandModule
	_, err := cm.Main("sh", "-c", "exit 1")
	if err == nil {
		t.Fatal("command should failed, but not")
	}
}
