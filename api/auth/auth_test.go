package auth_test

import (
	"testing"

	"github.com/pj-aias/matching-app-server/auth"
)

func TestPasswordValidation(t *testing.T) {
	rawPassword := "abc"

	hashed, err := auth.GeneratePasswordHash(rawPassword)
	if err != nil {
		panic(err)
	}

	got := auth.ValidatePassword([]byte(hashed), rawPassword)
	var wanted error = nil

	if got != wanted {
		t.Errorf("auth.ValidatePassword() failed; err = %v", got)
	}
}
