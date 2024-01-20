package services

import (
	"context"
	"fmt"
	"time"
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

func (p *Publication) DeletePublication(id string)  error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := fmt.Sprintf("DELETE FROM publications WHERE id = '%s'", id)
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
