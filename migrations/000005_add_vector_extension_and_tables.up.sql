-- Enable the pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Create a new table for embeddings
CREATE TABLE publication_embeddings (
    publication_id BIGSERIAL PRIMARY KEY REFERENCES publications(id),
    embedding vector(768),
    publication_type varchar(10),
    genres text[]
);

-- Create an index on the embedding column for faster similarity searches
CREATE INDEX ON publication_embeddings USING ivfflat (embedding vector_cosine_ops);