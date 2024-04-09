package services

import (
	"context"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserPublicationFollows struct {
	ID            uint      `json:"id"`             //"id" bigserial PRIMARY KEY,
	UserID        uint      `json:"user_id"`        //"user_id" bigint NOT NULL,
	PublicationID uint      `json:"publication_id"` //"publication_id" bigint NOT NULL,
	ChapterID     uint      `json:"chapter_id"`     //"chapter_id" bigint,
	CreatedAt     time.Time `json:"created_at"`     //"created_at" timestamptz DEFAULT (now()),
	UpdatedAt     time.Time `json:"updated_at"`     //"updated_at" timestamptz DEFAULT (now())
}

func (f *UserPublicationFollows) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.PublicationID, validation.Required),
		validation.Field(&f.ChapterID, is.Int),
	)
}

func (f *UserPublicationFollows) GetAllUserPublicationFollows(userID uint) ([]*UserPublicationFollows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
		SELECT f.id, f.user_id, f.publication_id, f.chapter_id 
		FROM user_publication_follows f 
		WHERE f.user_id = $1
	`

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*UserPublicationFollows
	for rows.Next() {
		var userPublicationFollows UserPublicationFollows
		err := rows.Scan(
			&userPublicationFollows.ID,
			&userPublicationFollows.UserID,
			&userPublicationFollows.PublicationID,
			&userPublicationFollows.ChapterID,
			// &userPublicationFollows.CreatedAt,
			// &userPublicationFollows.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, &userPublicationFollows)
	}

	fmt.Println(query)
	return data, nil
}

func (f *UserPublicationFollows) GetUserPublicationFollowsByPublicationID(ctx context.Context, publicationID string) (*UserPublicationFollows, error) {
	userPublicationFollows := &UserPublicationFollows{}

	query := `
		SELECT f.id, f.user_id, f.publication_id, f.chapter_id
		FROM user_publication_follows
		WHERE f.publication_id = $1
	`

	err := db.QueryRowContext(ctx, query, publicationID).Scan(
		&userPublicationFollows.ID,
		&userPublicationFollows.UserID,
		&userPublicationFollows.PublicationID,
		&userPublicationFollows.ChapterID,
	)
	if err != nil {
		return nil, err
	}

	return userPublicationFollows, nil
}

func (f *UserPublicationFollows) CreateUserPublicationFollows(ctx context.Context, update UserPublicationFollows) (*UserPublicationFollows, error) {
	query := `
		INSERT INTO user_publication_follows (user_id, publication_id, chapter_id, created_at, updated_at) 
		VALUES ($1, $2, $3, &4, &5)
		RETURNING id, created_at, updated_at
	`
	err := db.QueryRowContext(
		ctx, 
		query,
		update.ID,
		update.PublicationID,
		update.ChapterID,
		time.Now(),
		time.Now(),
	).Scan(
		&update.ID,
		&update.CreatedAt,
		&update.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &update, nil
}
