package group

import (
	"strings"
	"testing"

	"github.com/hnakamur/commango/os/unix/user"
)

func TestLookup(t *testing.T) {
	user, err := user.Current(nil)
	if err != nil {
		t.Fatal(err)
	}

	groups, err := AllGroups()
	if err != nil {
		t.Fatal(err)
	}

	group, err := LookupId(user.Gid, groups)
	if err != nil {
		t.Fatal(err)
	}

	group2, err := Lookup(group.Name, groups)
	if err != nil {
		t.Fatal(err)
	}

	if group2.Gid != group.Gid {
		t.Fatalf("excepted gid %s, got %s", group.Gid, group2.Gid)
	}
}

func TestSupplementaryGroups(t *testing.T) {
	_, err := SupplementaryGroups("daemon", nil)
	if err != nil {
		t.Fatal(err)
	}
}
