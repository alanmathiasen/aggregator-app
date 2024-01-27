package services

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Publication struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Rating float64 `json:"rating"`
	Image string `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Publication) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Title, validation.Required, validation.Length(3, 50)),
		validation.Field(&p.Description, validation.Required, validation.Length(3, 300)),
		validation.Field(&p.Rating, validation.Min(0.0), validation.Max(5.0)),
		validation.Field(&p.Image, validation.Required, is.URL),
	)
}

func (p *Publication) GetAllPublications() ([]*Publication, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, title, description, rating, image, created_at, updated_at FROM publications`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	
	var publications []*Publication
	for rows.Next() {
		var publication Publication
		err := rows.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Description,
			&publication.Rating,
			&publication.Image,
			&publication.CreatedAt,
			&publication.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		publications = append(publications, &publication)
	}

	return publications, nil
}

func (p *Publication) GetPublicationById(id string) (*Publication, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, title, description, rating, image, created_at, updated_at 
		FROM publications 
		WHERE id = $1
	`
	publication := &Publication{}
	
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&publication.ID,
		&publication.Title,
		&publication.Description,
		&publication.Rating,
		&publication.Image,
		&publication.CreatedAt,
		&publication.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return publication, nil
}

func (p *Publication) CreatePublication(publication Publication) (*Publication, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO publications (title, description, rating, image, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at
		`

	err := db.QueryRowContext(ctx, 
		query,
		publication.Title,
		publication.Description,
		publication.Rating,
		publication.Image,
		time.Now(),
		time.Now(),
	).Scan(&publication.ID, &publication.CreatedAt, &publication.UpdatedAt)
	
	if err != nil{
		return nil, err
	}
	
	return &publication, nil
}

func (p *Publication) UpdatePublication(id string, update Publication) (*Publication, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()



	query := `
		UPDATE publications
		SET
			title = $1,
			description = $2,
			rating = $3,
			image = $4,
			updated_at = $5
		WHERE id = $6
		RETURNING id, created_at, updated_at
	`
	err := db.QueryRowContext(
		ctx,
		query,
		update.Title,
		update.Description,
		update.Rating,
		update.Image,
		time.Now(),
		id,
	).Scan(&update.ID, &update.CreatedAt, &update.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &update, nil
}

func (p *Publication) DeletePublication(id string)  error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM publications WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
