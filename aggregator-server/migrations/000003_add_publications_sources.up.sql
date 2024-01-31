CREATE TABLE IF NOT EXISTS publication_sources (
    "id" uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "publication_id" uuid NOT NULL,
    "link" varchar NOT NULL,
    "source" varchar NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY ("publication_id") REFERENCES publications ("id")
);