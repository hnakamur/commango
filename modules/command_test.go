package modules

import (
	"testing"

	"github.com/hnakamur/commango/jsonutil"
)

func TestCommandNoArg(t *testing.T) {
	result, err := Command("hostname")
	if err != nil {
		t.Fatal(err)
	}

	_, err = jsonutil.Encode(result)
	if err != nil {
		t.Fatal(err)
	}
}

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
	result, err := Command("sh", "-c", "exit 1")
	if err == nil {
		t.Fatal("expected err is not nil")
	}
	if result.Rc != 1 {
		t.Fatal("expected Rc is 1, got %d", result.Rc)
	}
}
