package osutil_test

import (
	"io/ioutil"
	"os"
	"testing"
	"github.com/hnakamur/commango/os/osutil"
)

func TestExists(t *testing.T) {
	f, err := ioutil.TempFile("", "output-")
	if err != nil {
		t.Fatal(err)
	}

	if !osutil.Exists(f.Name()) {
		t.Fatal("Must exists")
	}

	err = os.Remove(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	if osutil.Exists(f.Name()) {
		t.Fatal("Must not exist")
	}
}
