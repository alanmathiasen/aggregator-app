package services

import (
	"context"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserPublicationFollows struct {
	ID            uint64    `json:"id"`
	UserID        uint      `json:"user_id"`
	PublicationID uint      `json:"publication_id"`
	ChapterID     uint      `json:"chapter_id"`
	Status        string      `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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

func (f *UserPublicationFollows) UpsertUserPublicationFollows(ctx context.Context, upf UserPublicationFollows) (*UserPublicationFollows, error) {
	query := `
		INSERT INTO user_publication_follows (user_id, publication_id, chapter_id, status, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, publication_id)
		DO UPDATE SET 
			chapter_id = EXCLUDED.chapter_id,
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at
		RETURNING id, created_at, updated_at
	`
	err := db.QueryRowContext(
		ctx,
		query,
		upf.UserID,
		upf.PublicationID,
		upf.ChapterID,
		upf.Status,
		time.Now(),
		time.Now(),
	).Scan(
		&upf.ID,
		&upf.CreatedAt,
		&upf.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &upf, nil
}

func (upf *UserPublicationFollows) DeleteUserPublicationFollow(ctx context.Context, publicationID uint, userID uint) error {
	query := `
		DELETE FROM user_publication_follows upf
		WHERE upf.publication_id = $1 AND upf.user_id = $2
	`
	_, err := db.ExecContext(ctx, query, publicationID, userID)
	if err != nil  {
		return err
	}
	return nil
}