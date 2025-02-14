package service

import (
	"golang.org/x/crypto/bcrypt"

	internalErrors "merch-shop/internal/errors"
	"merch-shop/pkg/logger"
)

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Service.hashPassword: ", "message", "error hashing password", "error", err)
		return "", internalErrors.ErrHashingPassword
	}
	return string(hash), nil
}

func comparePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
