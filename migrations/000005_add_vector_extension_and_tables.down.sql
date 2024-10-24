-- Down migration
-- Drop the index
DROP INDEX IF EXISTS publication_embeddings_embedding_idx;

-- Drop the table
DROP TABLE IF EXISTS publication_embeddings;

-- Disable the pgvector extension
DROP EXTENSION IF EXISTS vector;