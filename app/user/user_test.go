package user

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestUserBeforeCreateHook(t *testing.T) {
	u := User{Password: "rahasia"}
	u.BeforeCreate()

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("rahasia")); err != nil {
		t.Error("Password not match")
	}
}
