package services

import (
	"context"
	"time"
)

type User struct {
	ID string `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	
	err := db.QueryRowContext(
		ctx, 
		"SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1", 
		email,
		).Scan(
			&user.ID, 
			&user.Email, 
			&user.Password, 
			&user.CreatedAt, 
			&user.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *User) CreateUser(ctx context.Context, email string, hashedPassword string) (*User, error) {
	query := "INSERT INTO users (email, password, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, email, created_at, updated_at"
	err := db.QueryRowContext(
		ctx, 
		query, 
		email, 
		hashedPassword, 
		time.Now(), 
		time.Now(),
	).Scan(
		&u.ID, 
		&u.Email, 
		&u.CreatedAt, 
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return u, nil

}