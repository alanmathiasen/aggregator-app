package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sort"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgtype"
)

type Source struct {
	Link   string `json:"link"`
	Source string `json:"source"`
}

type Publication struct {
	ID                    int         `json:"id"`
	Title                 string      `json:"title"`
	Sinopsis              string      `json:"description"`
	Finished              bool        `json:"finished"`
	ReleaseDate           pgtype.Date `json:"firstAired"`
	TypeID                int64
	AuthorID              int64
	StudioID              int64
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	Status                *string    `json:"status"`
	LastReadChapterId     *int       `json:"chapter_id"`
	LastReadChapterNumber *string    `json:"last_read_chapter_number"`
	LastChapterReadAt     *time.Time `json:"last_chapter_read_at"`
	// IsFollowed              bool       `json:"is_followed"`
	// Image                   string     `json:"image"`
	// Sources                 []*Source  `json:"sources"`
	// PublicationFollowStatus *string    `json:"status"`
	// LastReadChapterId       *int       `json:"chapter_id"`
	// LastReadChapterNumber *string `json:"last_read_chapter_number"`
	// LastChapterReadAt       *time.Time `json:"last_chapter_read_at"`
	// Chapters                []*Chapter `json:"chapters"`
}

func (p *Publication) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Title, validation.Required, validation.Length(3, 50)),
		validation.Field(&p.Sinopsis, validation.Required, validation.Length(3, 300)),
		// validation.Field(&p.Image, validation.Required, is.URL),
	)
}

func (p *Publication) fetchPublications(ctx context.Context, userID uint, id string) ([]*Publication, error) {
	baseQuery := `
        SELECT 
			p.id, 
			p.title, 
			p.sinopsis, 
			p.finished,
			p.release_date,
			p.type_id,
			p.created_at,
			p.updated_at,
			ubs.name AS status, 
			ub.last_chapter_interacted_id,
			c2.number AS last_read_chapter_number,
			ub.last_chapter_interacted_at
		FROM publications p
		LEFT JOIN user_bookmarks ub ON p.id = ub.publication_id AND ub.user_id = $1
		LEFT JOIN chapters c ON p.id = c.publication_id 
		LEFT JOIN chapters c2 ON ub.last_chapter_interacted_id = c2.id
		LEFT JOIN user_bookmark_status ubs ON ub.status_id = ubs.id
	`
	// 	baseQuery := `
	// 	SELECT
	// 	p.id,
	// 	p.title,
	// 	p.sinopsis,
	// 	p.image,
	// 	CASE
	// 		WHEN ub.publication_id IS NULL THEN false
	// 		WHEN ub.status = 'deleted' THEN false
	// 		ELSE true
	// 	END AS is_followed,
	// 	ub.status,
	// 	ub.chapter_id AS last_read_chapter_id,
	// 	c2.number AS last_read_chapter_number,
	// 	ub.updated_at AS last_chapter_read_at,
	// 	json_agg(c) FILTER (WHERE c IS NOT NULL) AS chapters
	// 	FROM publications p
	// 	LEFT JOIN user_bookmarks ub ON p.id = ub.publication_id AND ub.user_id = $1
	// 	LEFT JOIN chapters c ON p.id = c.publication_id
	// 	LEFT JOIN chapters c2 ON ub.chapter_id = c2.id
	// `
	var rows *sql.Rows
	var err error

	if id != "" {
		query := baseQuery + " WHERE p.id = $2 GROUP BY p.id, ub.publication_id, ubs.name, ub.last_chapter_interacted_id, ub.last_chapter_interacted_at, c2.number"
		rows, err = db.QueryContext(ctx, query, userID, id)
	} else {
		query := baseQuery + " GROUP BY p.id, ub.publication_id, ubs.name, ub.last_chapter_interacted_id, ub.last_chapter_interacted_at, c2.number"
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
			&publication.Sinopsis,
			&publication.Finished,
			&publication.ReleaseDate,
			&publication.TypeID,
			&publication.CreatedAt,
			&publication.UpdatedAt,
			&publication.Status,
			&publication.LastReadChapterId,
			&publication.LastReadChapterNumber,
			&publication.LastChapterReadAt,
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
			// publication.Chapters = chapters
		} else {
			// publication.Chapters = make([]*Chapter, 0)
		}

		publications = append(publications, &publication)
	}

	return publications, nil
}

type GetAllPublicationsOptions struct {
	UserID uint
	Sort   string
	Page   int
	Limit  int
}

func (p *Publication) GetAllPublications(ctx context.Context, opts GetAllPublicationsOptions) ([]*Publication, error) {
	query := `
		SELECT 
			p.id, 
			p.title, 
			p.sinopsis, 
			p.finished,
			p.release_date,
			p.type_id
		FROM publications p
		ORDER BY p.release_date DESC
		LIMIT $1 OFFSET $2`

	var rows *sql.Rows
	var err error

	var limit int
	if opts.Limit == 0 {
		limit = 20
	} else {
		limit = opts.Limit
	}

	rows, err = db.QueryContext(ctx, query, limit, opts.Page*limit)
	if err != nil {
		return nil, err
	}

	var publications []*Publication
	for rows.Next() {
		var publication Publication
		err := rows.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Sinopsis,
			&publication.Finished,
			&publication.ReleaseDate,
			&publication.TypeID,
		)
		if err != nil {
			return nil, err
		}
		publications = append(publications, &publication)
	}
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

func (p *Publication) GetAllPublicationsByUser(ctx context.Context, userID uint) ([]*Publication, error) {
	return p.fetchPublications(ctx, userID, "")
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
		publication.Sinopsis,
		// publication.Image,
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
		update.Sinopsis,
		// update.Image,
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
