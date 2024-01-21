package tools

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) (hashedPwd string, err error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	hashedPwd = string(hashed)
	return
}

func IsValidPassword(hashed string, pwd string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashed),
		[]byte(pwd),
	)

	return err == nil
}
