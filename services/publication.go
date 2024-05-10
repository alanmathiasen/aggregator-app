package services

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Source struct {
	Link   string `json:"link"`
	Source string `json:"source"`
}

type Publication struct {
	ID                      int       `json:"id"`
	Title                   string    `json:"title"`
	Description             string    `json:"description"`
	Image                   string    `json:"image"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	Sources                 []*Source `json:"sources"`
	IsFollowed              bool      `json:"is_followed"`
	PublicationFollowStatus *string   `json:"status"`
	LastReadChapterId       *int      `json:"chapter_id"`
	LastReadChapterNumber   *string   `json:"chapter_number"`
	LastChapterReadAt       *string   `json:"chapter_followed"`
	ChapterNumbers          []string  `json:"chapter_numbers"`
}

func (p *Publication) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Title, validation.Required, validation.Length(3, 50)),
		validation.Field(&p.Description, validation.Required, validation.Length(3, 300)),
		validation.Field(&p.Image, validation.Required, is.URL),
	)
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
			CASE 
				WHEN upf.publication_id IS NULL THEN false 
				WHEN upf.status = 'deleted' THEN false
				ELSE true 
			END AS is_followed,
			upf.status, 
			upf.chapter_id AS last_chapter_read_id,
			c.number AS last_read_chapter_number,
			upf.updated_at AS last_chapter_read_at,
			ARRAY_AGG(c.number) AS chapter_numbers 
			FROM publications p
			LEFT JOIN user_publication_follows upf ON p.id = upf.publication_id AND upf.user_id = $1
			LEFT JOIN chapters c ON upf.chapter_id = c.id
			LEFT JOIN chapters c2 ON p.id = c2.publication_id
			WHERE p.id = $2
			GROUP BY p.id, upf.publication_id, upf.status, upf.chapter_id, upf.updated_at, c.number
	`

	// 		ps.link,
	// 		ps.source
	// 	FROM publications p
	// 	LEFT JOIN user_publication_follows upf ON p.id = upf.publication_id AND upf.user_id = 1
	// 	LEFT JOIN chapters c ON upf.chapter_id = c.id
	// 	LEFT JOIN publication_sources ps ON p.id = ps.publication_id
	// 	WHERE p.id = $1
	publication := &Publication{}
	var chapterNumbers string

	err := db.QueryRowContext(ctx, query, userID, id).Scan(
		&publication.ID,
		&publication.Title,
		&publication.Description,
		&publication.Image,
		&publication.IsFollowed,
		&publication.PublicationFollowStatus,
		&publication.LastReadChapterId,
		&publication.LastReadChapterNumber,
		&publication.LastChapterReadAt,
		&chapterNumbers,
		// &publication.CreatedAt,
		// &publication.UpdatedAt,
		// &source.Link,
		// &source.Source,
	)

	if err != nil {
		return nil, err
	}
	numbers := strings.Split(strings.Trim(chapterNumbers, "{}"), ",")
	chapters := make([]int, len(numbers))
	for i, number := range numbers {
		chapters[i], _ = strconv.Atoi(number)
	}
	sort.Ints(chapters)
	sortedChapters := make([]string, len(chapters))
	for i, chapter := range chapters {
		sortedChapters[i] = fmt.Sprint(chapter)
	}
	publication.ChapterNumbers = sortedChapters

	// publication.Sources = append(publication.Sources, &source)
	// }
	return publication, nil
}

func (p *Publication) GetAllPublications(ctx context.Context, userID uint) ([]*Publication, error) {
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
	upf.chapter_id AS last_read_chapter_id,
	c2.number AS last_read_chapter_number,
  	upf.updated_at AS last_chapter_read_at, 
	ARRAY_AGG(c.number) AS chapter_numbers
	FROM publications p
	LEFT JOIN user_publication_follows upf ON p.id = upf.publication_id AND upf.user_id = $1
	LEFT JOIN chapters c ON p.id = c.publication_id 
	LEFT JOIN chapters c2 ON upf.chapter_id = c2.id
	GROUP BY p.id, upf.publication_id, upf.status, upf.chapter_id, upf.updated_at, c2.number
`

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	var publications []*Publication
	for rows.Next() {
		var publication Publication
		var chapterNumbers string

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
			&chapterNumbers,
			//&publication.CreatedAt,
			//&publication.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		numbers := strings.Split(strings.Trim(chapterNumbers, "{}"), ",")
		chapters := make([]int, len(numbers))
		for i, number := range numbers {
			chapters[i], _ = strconv.Atoi(number)
		}
		sort.Ints(chapters)
		sortedChapters := make([]string, len(chapters))
		for i, chapter := range chapters {
			sortedChapters[i] = fmt.Sprint(chapter)
		}
		publication.ChapterNumbers = sortedChapters
		publications = append(publications, &publication)
	}

	// pubs, err := json.MarshalIndent(publications, "", " ")
	// if err != nil {
	// 	return nil, err
	// }

	// fmt.Println(string(pubs))

	return publications, nil
}

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
