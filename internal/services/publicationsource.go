package services

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type PublicationSource struct {
	ID            string    `json:"id"`
	PublicationID string    `json:"publication_id"`
	Link          string    `json:"link"`
	Source        string    `json:"source"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (ps *PublicationSource) Validate() error {
	return validation.ValidateStruct(
		ps,
		validation.Field(&ps.PublicationID, validation.Required, is.UUIDv4),
		validation.Field(&ps.Link, validation.Required, is.URL),
	)
}

func (ps *PublicationSource) GetAllLinksByPublicationID(publicationID string) ([]*PublicationSource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, publication_id, link, created_at, updated_at FROM publication_links WHERE publication_id = $1`
	rows, err := db.QueryContext(ctx, query, publicationID)
	if err != nil {
		return nil, err
	}

	var publicationSources []*PublicationSource
	for rows.Next() {
		var publicationSource PublicationSource
		err := rows.Scan(
			&publicationSource.ID,
			&publicationSource.PublicationID,
			&publicationSource.Link,
			&publicationSource.CreatedAt,
			&publicationSource.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		publicationSources = append(publicationSources, &publicationSource)
	}

	return publicationSources, nil
}

func (ps *PublicationSource) CreatePublicationSource(publicationSource PublicationSource, publicationID string) (*PublicationSource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
	INSERT INTO publication_links (publication_id, link, created_at, updated_at)
	VALUES ($1, $2, $3, $4) 
	RETURNING id, created_at, updated_at
	`
	err := db.QueryRowContext(
		ctx,
		query,
		publicationID,
		publicationSource.Link,
		time.Now(),
		time.Now(),
	).Scan(
		&publicationSource.ID,
		&publicationSource.CreatedAt,
		&publicationSource.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &publicationSource, nil
}
