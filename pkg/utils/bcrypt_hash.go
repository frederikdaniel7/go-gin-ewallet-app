package utils

import "golang.org/x/crypto/bcrypt"

type Crypto interface {
	HashPassword(password string, cost int) ([]byte, error)
	CheckPassword(password string, hash []byte) (bool, error)
}

func HashPassword(pwd string, cost int) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), cost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func CheckPassword(pwd string, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(pwd))
	if err != nil {
		return false, err
	}
	return true, nil
}
