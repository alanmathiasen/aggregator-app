ALTER TABLE chapters ADD CONSTRAINT unique_publication_season_number UNIQUE (publication_id, season_number, number);