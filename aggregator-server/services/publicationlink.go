package services

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type PublicationLink struct {
	ID string `json:"id"`
	PublicationID string `json:"publication_id"`
	Link string `json:"link"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (pl *PublicationLink) Validate() error {
	return validation.ValidateStruct(
		pl,
		validation.Field(&pl.PublicationID, validation.Required, is.UUIDv4),
		validation.Field(&pl.Link, validation.Required, is.URL),
	)
}

func (pl *PublicationLink) GetAllLinksByPublicationID(publicationID string) ([]*PublicationLink, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	
	query := `SELECT id, publication_id, link, created_at, updated_at FROM publication_links WHERE publication_id = $1`
	rows, err := db.QueryContext(ctx, query, publicationID)
	if err != nil {
		return nil, err
	}

	var publicationLinks []*PublicationLink
	for rows.Next() {
		var publicationLink PublicationLink
		err := rows.Scan(
			&publicationLink.ID,
			&publicationLink.PublicationID,
			&publicationLink.Link,
			&publicationLink.CreatedAt,
			&publicationLink.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		publicationLinks = append(publicationLinks, &publicationLink)
	} 
	
	return publicationLinks, nil
}

func (pl *PublicationLink) CreatePublicationLink(publicationLink PublicationLink, publicationID string) (*PublicationLink, error) {
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
		publicationLink.Link,
		time.Now(),
		time.Now(),
	).Scan(
		&publicationLink.ID, 
		&publicationLink.CreatedAt, 
		&publicationLink.UpdatedAt,
	)

	if err != nil{
		return nil, err
	}

	return &publicationLink, nil
}