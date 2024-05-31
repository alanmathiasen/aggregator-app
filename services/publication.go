package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Source struct {
	Link   string `json:"link"`
	Source string `json:"source"`
}

type Publication struct {
	ID                      int        `json:"id"`
	Title                   string     `json:"title"`
	Description             string     `json:"description"`
	Image                   string     `json:"image"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
	Sources                 []*Source  `json:"sources"`
	IsFollowed              bool       `json:"is_followed"`
	PublicationFollowStatus *string    `json:"status"`
	LastReadChapterId       *int       `json:"chapter_id"`
	LastReadChapterNumber   *string    `json:"last_read_chapter_number"`
	LastChapterReadAt       *time.Time `json:"last_chapter_read_at"`
	Chapters                []*Chapter `json:"chapters"`
}

func (p *Publication) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Title, validation.Required, validation.Length(3, 50)),
		validation.Field(&p.Description, validation.Required, validation.Length(3, 300)),
		validation.Field(&p.Image, validation.Required, is.URL),
	)
}

func (p *Publication) fetchPublications(ctx context.Context, userID uint, id string) ([]*Publication, error) {

	baseQuery := `
        SELECT 
        p.id, 
        p.title, 
        p.description, 
        p.image, 
        CASE 
            WHEN upf.publication_id IS NULL THEN false 
            WHEN upf.status = 'deleted' THEN false
            ELSE true 
        END AS is_followed,
        upf.status, 
        upf.chapter_id AS last_read_chapter_id,
        c2.number AS last_read_chapter_number,
        upf.updated_at AS last_chapter_read_at, 
		json_agg(c) FILTER (WHERE c IS NOT NULL) AS chapters
        FROM publications p
        LEFT JOIN user_publication_follows upf ON p.id = upf.publication_id AND upf.user_id = $1
        LEFT JOIN chapters c ON p.id = c.publication_id 
        LEFT JOIN chapters c2 ON upf.chapter_id = c2.id
    `
	//COALESCE(json_agg(c)) AS chapters

	var rows *sql.Rows
	var err error

	if id != "" {
		query := baseQuery + " WHERE p.id = $2 GROUP BY p.id, upf.publication_id, upf.status, upf.chapter_id, upf.updated_at, c2.number"
		rows, err = db.QueryContext(ctx, query, userID, id)
	} else {
		query := baseQuery + " GROUP BY p.id, upf.publication_id, upf.status, upf.chapter_id, upf.updated_at, c2.number"
		rows, err = db.QueryContext(ctx, query, userID)
	}

	if err != nil {
		return nil, err
	}

	var publications []*Publication
	for rows.Next() {
		var publication Publication
		var chaptersJSON sql.NullString

		err := rows.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Description,
			&publication.Image,
			&publication.IsFollowed,
			&publication.PublicationFollowStatus,
			&publication.LastReadChapterId,
			&publication.LastReadChapterNumber,
			&publication.LastChapterReadAt,
			&chaptersJSON,
		)
		if err != nil {
			return nil, err
		}
		if chaptersJSON.Valid && chaptersJSON.String != "" {
			var chapters []*Chapter
			if err := json.Unmarshal([]byte(chaptersJSON.String), &chapters); err != nil {
				return nil, err
			}
			sort.Slice(chapters, func(i, j int) bool {
				return chapters[i].Number < chapters[j].Number
			})
			publication.Chapters = chapters
		} else {
			publication.Chapters = make([]*Chapter, 0)
		}
		// var chapters []*Chapter
		// if err := json.Unmarshal([]byte(chaptersJSON), &chapters); err != nil {
		// 	return nil, err
		// }
		// fmt.Println(string(chaptersJSON))
		// sort.Slice(chapters, func(i, j int) bool {
		// 	return chapters[i].Number > chapters[j].Number
		// })
		// publication.Chapters = chapters
		publications = append(publications, &publication)
	}

	data, err := json.MarshalIndent(publications, "", "  ")
	fmt.Println(string(data))
	return publications, nil
}

func (p *Publication) GetPublicationById(ctx context.Context, id string, userID uint) (*Publication, error) {
	publications, err := p.fetchPublications(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	if len(publications) == 0 {
		return nil, errors.New("no publication found")
	}

	return publications[0], nil
}

func (p *Publication) GetAllPublications(ctx context.Context, userID uint) ([]*Publication, error) {
	return p.fetchPublications(ctx, userID, "")
}

// func (p *Publication) GetPublicationById(id string, userID uint) (*Publication, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
// 	defer cancel()

// 	query := `
// 		SELECT
// 			p.id,
// 			p.title,
// 			p.description,
// 			p.image,
// 			CASE
// 				WHEN upf.publication_id IS NULL THEN false
// 				WHEN upf.status = 'deleted' THEN false
// 				ELSE true
// 			END AS is_followed,
// 			upf.status,
// 			upf.chapter_id AS last_chapter_read_id,
// 			c2.number AS last_read_chapter_number,
// 			upf.updated_at AS last_chapter_read_at,
// 			COALESCE(json_agg(c)) AS chapters
// 			FROM publications p
// 			LEFT JOIN user_publication_follows upf ON p.id = upf.publication_id AND upf.user_id = $1
// 			LEFT JOIN chapters c ON p.id = c.publication_id
// 			LEFT JOIN chapters c2 ON upf.chapter_id = c2.id
// 			WHERE p.id = $2
// 			GROUP BY p.id, upf.publication_id, upf.status, upf.chapter_id, upf.updated_at, c2.number
// 	`
// 	publication := &Publication{}
// 	var chaptersJSON string

// 	err := db.QueryRowContext(ctx, query, userID, id).Scan(
// 		&publication.ID,
// 		&publication.Title,
// 		&publication.Description,
// 		&publication.Image,
// 		&publication.IsFollowed,
// 		&publication.PublicationFollowStatus,
// 		&publication.LastReadChapterId,
// 		&publication.LastReadChapterNumber,
// 		&publication.LastChapterReadAt,
// 		&chaptersJSON,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}
// 	var chapters []Chapter
// 	if err := json.Unmarshal([]byte(chaptersJSON), &chapters); err != nil {
// 		return nil, err
// 	}
// 	sort.Slice(chapters, func(i, j int) bool {
// 		return chapters[i].Number > chapters[j].Number
// 	})
// 	publication.Chapters = chapters

// 	return publication, nil
// }

// func (p *Publication) GetAllPublications(ctx context.Context, userID uint) ([]*Publication, error) {
// 	query := `
// 		SELECT
// 		p.id,
// 		p.title,
// 		p.description,
// 		p.image,
// 		CASE
// 			WHEN upf.publication_id IS NULL THEN false
// 			WHEN upf.status = 'deleted' THEN false
// 			ELSE true
// 		END AS is_followed,
// 		upf.status,
// 		upf.chapter_id AS last_read_chapter_id,
// 		c2.number AS last_read_chapter_number,
// 		upf.updated_at AS last_chapter_read_at,
// 		COALESCE(json_agg(c)) AS chapters
// 		FROM publications p
// 		LEFT JOIN user_publication_follows upf ON p.id = upf.publication_id AND upf.user_id = $1
// 		LEFT JOIN chapters c ON p.id = c.publication_id
// 		LEFT JOIN chapters c2 ON upf.chapter_id = c2.id
// 		GROUP BY p.id, upf.publication_id, upf.status, upf.chapter_id, upf.updated_at, c2.number
// 	`

// 	rows, err := db.QueryContext(ctx, query, userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var publications []*Publication
// 	for rows.Next() {
// 		var publication Publication
// 		var chaptersJSON string

// 		err := rows.Scan(
// 			&publication.ID,
// 			&publication.Title,
// 			&publication.Description,
// 			&publication.Image,
// 			&publication.IsFollowed,
// 			&publication.PublicationFollowStatus,
// 			&publication.LastReadChapterId,
// 			&publication.LastReadChapterNumber,
// 			&publication.LastChapterReadAt,
// 			&chaptersJSON,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		var chapters []Chapter
// 		if err := json.Unmarshal([]byte(chaptersJSON), &chapters); err != nil {
// 			return nil, err
// 		}
// 		sort.Slice(chapters, func(i, j int) bool {
// 			return chapters[i].Number > chapters[j].Number
// 		})
// 		publication.Chapters = chapters
// 		publications = append(publications, &publication)
// 	}

// 	return publications, nil
// }

func (p *Publication) GetUserPublicationsHTML(ctx context.Context, userID uint) ([]*Publication, error) {
	query := `
			SELECT 
			p.id, 
			p.title, 
			p.description, 
			p.image, 
			CASE 
				WHEN upf.publication_id IS NULL THEN false 
				WHEN upf.status = 'deleted' THEN false
				ELSE true 
			END AS is_followed,
			upf.status, 
			upf.chapter_id AS last_chapter_read,
			c.id AS chapter_id,
			c.number AS chapter_number
			ARRAY_AGG(c2.number) AS chapter_numbers
			FROM publications p
			LEFT JOIN user_publication_follows upf ON p.id = upf.publication_id AND upf.user_id = $1
			LEFT JOIN chapters c ON upf.chapter_id = c.id 
			LEFT JOIN chapters c2 ON p.id = c2.publication_id
			GROUP BY p.id, upf.publication_id, upf.status, upf.chapter_id, upf.updated_at, c2.number
		`
	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	var publications []*Publication
	for rows.Next() {
		publication := &Publication{}
		err := rows.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Image,
			&publication.IsFollowed,
			&publication.PublicationFollowStatus,
			&publication.LastChapterReadAt,
			&publication.LastReadChapterId,
			&publication.LastReadChapterNumber,
		)

		if err != nil {
			return nil, err
		}
		fmt.Println("aca en el service bro y vos?", publication.PublicationFollowStatus)
		publications = append(publications, publication)
	}
	return publications, nil
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
