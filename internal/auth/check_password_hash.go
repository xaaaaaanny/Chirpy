package auth

import "golang.org/x/crypto/bcrypt"

func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	
	return nil
}
