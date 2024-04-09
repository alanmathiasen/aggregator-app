package services

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Chapter struct {
	ID            string    `json:"id"`
	PublicationID string    `json:"publication_id"`
	Number        int       `json:"number"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Rating        float64   `json:"rating"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (c *Chapter) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.Number, validation.Min(0), validation.Max(5)),
		validation.Field(&c.Title, validation.Required, validation.Length(3, 50)),
		validation.Field(&c.Description, validation.Length(3, 300)),
		validation.Field(&c.Rating, validation.Min(0.0), validation.Max(5.0)),
	)
}

func (c *Chapter) GetAllChaptersByPublicationID(ctx context.Context, publicationID string) ([]*Chapter, error) {
	var chapters []*Chapter
	query := `SELECT id, publication_id, number, title, description, rating, created_at, updated_at FROM chapters WHERE publication_id = $1`
	rows, err := db.QueryContext(ctx, query, publicationID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var chapter Chapter
		err := rows.Scan(
			&chapter.ID,
			&chapter.PublicationID,
			&chapter.Number,
			&chapter.Title,
			&chapter.Description,
			&chapter.Rating,
			&chapter.CreatedAt,
			&chapter.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		chapters = append(chapters, &chapter)
	}

	return chapters, nil
}

func (c *Chapter) CreateChapterForPublication(ctx context.Context, publicationID string, data *Chapter) error {
	query := `
		INSERT INTO chapters (publication_id, number, title, description, rating, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, publication_id, created_at, updated_at	
	`
	err := db.QueryRowContext(
		ctx,
		query,
		publicationID,
		c.Number,
		c.Title,
		c.Description,
		c.Rating,
		time.Now(),
		time.Now(),
	).Scan(
		&data.ID,
		&data.PublicationID,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
