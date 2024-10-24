CREATE TABLE genres (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL
);

CREATE TABLE publication_genres (
    publication_id BIGSERIAL NOT NULL,
    genre_id BIGSERIAL NOT NULL,
    PRIMARY KEY (publication_id, genre_id),
    FOREIGN KEY (publication_id) REFERENCES publications(id) ON DELETE CASCADE,
    FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE
);


ALTER TABLE publications ADD COLUMN image_url VARCHAR(255);

