BEGIN;
ALTER TABLE publication_sources
DROP CONSTRAINT publication_sources_publication_id_fkey;

ALTER TABLE publication_sources
ADD CONSTRAINT publication_sources_publication_id_fkey
FOREIGN KEY (publication_id) REFERENCES publications(id) ON DELETE CASCADE;
COMMIT;