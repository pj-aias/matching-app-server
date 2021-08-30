package db_test

import (
	"os"
	"testing"
	"time"

	"github.com/pj-aias/matching-app-server/auth"
	"github.com/pj-aias/matching-app-server/db"
)

func TestPasswordHashCreate(t *testing.T) {
	rawPassword := "testpassword"

	user, err := db.LookupUser("passwordtest" + time.Now().String())
	if err != nil {
		user.Username = "passwordtest"
		user, err = db.AddUser(user)
		if err != nil {
			t.Errorf("failed to get / create test user: %v", err)
			return
		}
	}

	hashed, err := auth.GeneratePasswordHash(rawPassword)
	if err != nil {
		panic(err)
	}

	gotHash, err := db.AddPasswordHash(uint64(user.ID), hashed)
	if err != nil {
		t.Errorf("failed to insert hash: %v", err)
		return
	}

	got := gotHash.Hash
	wanted := hashed
	t.Logf("got\t: %v", got)
	t.Logf("wanted\t: %v", wanted)

	for i, g := range got {
		if g != wanted[i] {
			t.Errorf("invalid hash returned")
			t.Errorf("differs at [%v]: %v < %v", i, g, wanted[i])
			return
		}
	}

	if err := auth.ValidatePassword(got, rawPassword); err != nil {
		t.Errorf("invalid password got: %v", err)
		return
	}
	if err := auth.ValidatePassword(wanted, rawPassword); err != nil {
		t.Errorf("invalid password wanted: %v", err)
		return
	}
}

func TestMain(m *testing.M) {
	db.ConnectDB("test")
	code := m.Run()
	os.Exit(code)
}
