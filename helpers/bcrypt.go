package helpers

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// ? hashing password before saved in database
func HashPass(p string) string {
	password := []byte(p)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("Error hashing password:", err)
		return ""
	}

	return string(hash)
}

// ? comparing password has been hashing
func ComparePass(h, p []byte) bool {
	hash, pass := []byte(h), []byte(p)

	err := bcrypt.CompareHashAndPassword(hash, pass)

	return err == nil
}
