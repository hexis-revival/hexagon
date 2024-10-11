package common

import (
	"testing"
)

func TestPasswords(t *testing.T) {
	password := "password"
	hash, err := CreatePasswordHash(password)

	if err != nil {
		t.Error(err)
		return
	}

	if !CheckPassword(password, hash) {
		t.Error("password check failed")
		return
	}

	passwordHashed := GetSHA512Hash(password)

	if !CheckPasswordHashed(passwordHashed, hash) {
		t.Error("hashed password check failed")
		return
	}

	if len(passwordCache) != 1 {
		t.Error("password cache should have 1 entry")
		return
	}
}
