package helpers

import "golang.org/x/crypto/bcrypt"

func BcryptHashedPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return hash, err
}
