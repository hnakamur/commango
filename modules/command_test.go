package modules

import (
	"testing"

	"github.com/hnakamur/commango/jsonutil"
)

func TestCommand(t *testing.T) {
	var cm CommandModule
	rj, _ := cm.Main("sh", "-c", "exit 1")
	_, err := jsonutil.Encode(rj)
	if err != nil {
		t.Fatal(err)
	}
}
