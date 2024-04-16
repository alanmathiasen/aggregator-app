package services

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Source struct {
	Link   string `json:"link"`
	Source string `json:"source"`
}

type Publication struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Image         string    `json:"image"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Sources       []*Source `json:"sources"`
	IsFollowed    bool      `json:"is_followed"`
	Status        *string   `json:"status"`
	ChapterId     *int      `json:"chapter_id"`
	ChapterNumber *string   `json:"chapter_number"`
}

func (p *Publication) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Title, validation.Required, validation.Length(3, 50)),
		validation.Field(&p.Description, validation.Required, validation.Length(3, 300)),
		validation.Field(&p.Image, validation.Required, is.URL),
	)
}

func (p *Publication) GetAllPublications(userID uint) ([]*Publication, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// query := `SELECT id, title, description, image, created_at, updated_at FROM publications`
	// rows, err := db.QueryContext(ctx, query)
	// if err != nil {
	// 	return nil, err
	// }

	query := `
			SELECT 
			p.id, 
			p.title, 
			p.description, 
			p.image, 
			CASE WHEN upf.publication_id IS NULL THEN false ELSE true END AS is_followed,
			upf.status, 
			upf.chapter_id,
			c.number AS chapter_number
			FROM publications p
			LEFT JOIN user_publication_follows upf ON p.id = upf.publication_id AND upf.user_id = $1
			LEFT JOIN chapters c ON upf.chapter_id = c.id 
		`

	rows, err := db.QueryContext(ctx, query, userID)
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
			&publication.Image,
			&publication.IsFollowed,
			&publication.Status,
			&publication.ChapterId,
			&publication.ChapterNumber,
			//&publication.CreatedAt,
			//&publication.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		publications = append(publications, &publication)
	}

	return publications, nil
}

func (p *Publication) GetPublicationById(id string, userID uint) (*Publication, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			p.id, 
			p.title, 
			p.description, 
			p.image, 
			CASE WHEN upf.publication_id IS NULL THEN false ELSE true END AS is_followed,
			upf.status, 
			upf.chapter_id,
			c.number AS chapter_number
		FROM publications p
		LEFT JOIN user_publication_follows upf ON p.id = upf.publication_id AND upf.user_id = $2
		LEFT JOIN chapters c ON upf.chapter_id = c.id 
		WHERE p.id = $1
	`
	// SELECT
	// 		p.id,
	// 		p.title,
	// 		p.description,
	// 		p.image,
	// 		CASE WHEN upf.publication_id IS NULL THEN false ELSE true END AS is_followed,
	// 		upf.status,
	// 		upf.chapter_id,
	// 		c.number AS chapter_number,
	// 		ps.link,
	// 		ps.source
	// 	FROM publications p
	// 	LEFT JOIN user_publication_follows upf ON p.id = upf.publication_id AND upf.user_id = 1
	// 	LEFT JOIN chapters c ON upf.chapter_id = c.id
	// 	LEFT JOIN publication_sources ps ON p.id = ps.publication_id
	// 	WHERE p.id = $1
	publication := &Publication{}

	err := db.QueryRowContext(ctx, query, id, userID).Scan(
		&publication.ID,
		&publication.Title,
		&publication.Description,
		&publication.Image,
		&publication.IsFollowed,
		&publication.Status,
		&publication.ChapterId,
		&publication.ChapterNumber,
		// &publication.CreatedAt,
		// &publication.UpdatedAt,
		// &source.Link,
		// &source.Source,
	)
	// defer rows.Close()
	// for rows.Next() {
	// 	var source Source
	// 	err := rows.Scan(
	// 		&publication.ID,
	// 		&publication.Title,
	// 		&publication.Description,
	// 		&publication.Image,
	// 		&publication.IsFollowed,
	// 		&publication.Status,
	// 		&publication.ChapterId,
	// 		&publication.ChapterNumber,

	// 	)
	if err != nil {
		return nil, err
	}
	// publication.Sources = append(publication.Sources, &source)
	// }
	return publication, nil
}

func (p *Publication) CreatePublication(ctx context.Context, publication Publication) (*Publication, error) {
	query := `
		INSERT INTO publications (title, description, image, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at, updated_at
		`

	err := db.QueryRowContext(ctx,
		query,
		publication.Title,
		publication.Description,
		publication.Image,
		time.Now(),
		time.Now(),
	).Scan(&publication.ID, &publication.CreatedAt, &publication.UpdatedAt)
	if err != nil {
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
			image = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING id, created_at, updated_at
	`
	err := db.QueryRowContext(
		ctx,
		query,
		update.Title,
		update.Description,
		update.Image,
		time.Now(),
		id,
	).Scan(&update.ID, &update.CreatedAt, &update.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &update, nil
}

func (p *Publication) DeletePublication(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM publications WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
