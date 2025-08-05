package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func PrintHashedPassword(password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println("Hashed password:", string(hash))
}
