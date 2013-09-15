package user

import (
	"testing"
)

func TestCurrent(t *testing.T) {
	_, err := Current(nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLookup(t *testing.T) {
	users, err := AllUsers()
	if err != nil {
		t.Fatal(err)
	}

	user, err := Current(users)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Lookup(user.Username, users)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLookupId(t *testing.T) {
	users, err := AllUsers()
	if err != nil {
		t.Fatal(err)
	}

	user, err := Current(users)
	if err != nil {
		t.Fatal(err)
	}

	_, err = LookupId(user.Uid, users)
	if err != nil {
		t.Fatal(err)
	}
}
