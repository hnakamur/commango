package modules

import (
	"testing"

	"github.com/hnakamur/commango/jsonutil"
)

func TestCommand(t *testing.T) {
	var cm CommandModule
	//result, err := cm.Main("sh", "-c", "exit 1")
	result, err := cm.Main("uname", "-a")
	if err != nil {
		t.Fatal(err)
	}
	cr, ok := result.(CommandResult)
	if !ok {
		t.Fatal("CommandResult expected")
	}

	json, err := jsonutil.Encode(cr)
	t.Fatalf("json=%s.", json)

	//_, err := jsonutil.Encode(rj)
	if err != nil {
		t.Fatal(err)
	}
}
