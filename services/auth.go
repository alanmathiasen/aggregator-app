package services

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var user User

func AuthenticateUser(ctx context.Context, email string, password string) (*User, error) {
	userData, err := user.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !ComparePassword(userData.HashedPassword, password) {
		return nil, errors.New("incorrect password")
	}
	return userData, nil
}

func RegisterUser(ctx context.Context, email string, hashedPassword string) (*User, error) {
	newUser, err := user.CreateUser(ctx, email, hashedPassword)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ComparePassword(hashedPassword string, password string) bool {
	hashedPasswordBytes := []byte(hashedPassword)
	passwordBytes := []byte(password)
	return bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes) == nil
}
