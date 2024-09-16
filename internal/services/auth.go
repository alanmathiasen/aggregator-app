package services

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var user User

type AuthService struct{}

func (a *AuthService) AuthenticateUser(ctx context.Context, email string, password string) (*User, error) {
	userData, err := user.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("incorrect email")
	}

	if !a.ComparePassword(userData.HashedPassword, password) {
		return nil, errors.New("incorrect password")
	}

	return userData, nil
}

func (a *AuthService) RegisterUser(ctx context.Context, email string, hashedPassword string) (*User, error) {
	newUser, err := user.CreateUser(ctx, email, hashedPassword)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (a *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (a *AuthService) ComparePassword(hashedPassword string, password string) bool {
	hashedPasswordBytes := []byte(hashedPassword)
	passwordBytes := []byte(password)
	return bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes) == nil
}
