package yum

import (
	"testing"

	"github.com/hnakamur/commango/jsonutil"
)

func TestInstalled(t *testing.T) {
	result, err := Installed("kernel")
	if err != nil {
		t.Fatal(err)
	}

    _, err = jsonutil.Encode(result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotInstalled(t *testing.T) {
	result, err := Installed("no_such_package")
	if err != nil {
		t.Fatal(err)
	}

    _, err = jsonutil.Encode(result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInstallGroup(t *testing.T) {
	result, err := Install("@'Development tools'")
	if err != nil {
		t.Fatal(err)
	}

    _, err = jsonutil.Encode(result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInstall(t *testing.T) {
	result, err := Install("make")
	if err != nil {
		t.Fatal(err)
	}

    _, err = jsonutil.Encode(result)
	if err != nil {
		t.Fatal(err)
	}
}
