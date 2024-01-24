CREATE TABLE IF NOT EXISTS chapters (
    "id" uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "publication_id" uuid NOT NULL,
    "number" int NOT NULL,
    "title" varchar NOT NULL,
    "description" varchar NOT NULL,
    "rating" FLOAT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY ("publication_id") REFERENCES publications ("id")
);