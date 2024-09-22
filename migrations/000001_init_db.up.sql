CREATE TABLE "publications" (
  "id" bigserial PRIMARY KEY,
  "title" varchar,
  "sinopsis" text,
  "finished" bool,
  "release_date" timestamp,
  "type_id" bigint,
  "author_id" bigint,
  "studio_id" bigint,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "publication_types" (
  "id" bigserial PRIMARY KEY,
  "name" varchar,
  "created_at" timestamp
);

CREATE TABLE "publication_reviews" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint,
  "publication_id" bigint,
  "rating" float,
  "text" text,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "review_likes" (
  "id" bigserial PRIMARY KEY,
  "review_id" bigint,
  "user_id" bigint,
  "is_like" bool,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar,
  "role" varchar,
  "email" varchar,
  "hashed_password" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "chapters" (
  "id" bigserial PRIMARY KEY,
  "publication_id" bigint,
  "title" varchar,
  "number" varchar,
  "image" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "chapter_likes" (
  "id" bigserial PRIMARY KEY,
  "chapter_id" bigint,
  "user_id" bigint,
  "is_like" bool,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "chapter_comments" (
  "id" bigserial PRIMARY KEY,
  "chapter_id" bigint,
  "user_id" bigint,
  "text" text,
  "parent_id" bigint,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "chapter_comment_likes" (
  "id" bigserial PRIMARY KEY,
  "comment_id" bigint,
  "user_id" bigint,
  "is_like" bool,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "user_bookmarks" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint,
  "publication_id" bigint,
  "last_chapter_interacted_id" bigint,
  "status_id" bigint,
  "last_chapter_interacted_at" timestamp,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "user_bookmark_status" (
  "id" bigserial PRIMARY KEY,
  "name" varchar,
  "created_at" timestamp
);

CREATE TABLE "suggestions" (
  "id" bigserial PRIMARY KEY,
  "publication_1_id" bigint,
  "publication_2_id" bigint,
  "similarity_score" float,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "chapter_reviews" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint,
  "chapter_id" bigint,
  "rating" float,
  "text" text,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "chapter_review_likes" (
  "id" bigserial PRIMARY KEY,
  "review_id" bigint,
  "user_id" bigint,
  "is_like" bool,
  "created_at" timestamp,
  "updated_at" timestamp
);

-- Fix foreign key constraints
ALTER TABLE "publications" ADD FOREIGN KEY ("type_id") REFERENCES "publication_types" ("id");

ALTER TABLE "publication_reviews" ADD FOREIGN KEY ("publication_id") REFERENCES "publications" ("id");

ALTER TABLE "review_likes" ADD FOREIGN KEY ("review_id") REFERENCES "publication_reviews" ("id");

ALTER TABLE "review_likes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "publication_reviews" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "chapters" ADD FOREIGN KEY ("publication_id") REFERENCES "publications" ("id");

ALTER TABLE "chapter_likes" ADD FOREIGN KEY ("chapter_id") REFERENCES "chapters" ("id");

ALTER TABLE "chapter_likes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "chapter_comments" ADD FOREIGN KEY ("chapter_id") REFERENCES "chapters" ("id");

ALTER TABLE "chapter_comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "chapter_comments" ADD FOREIGN KEY ("parent_id") REFERENCES "chapter_comments" ("id");

ALTER TABLE "chapter_comment_likes" ADD FOREIGN KEY ("comment_id") REFERENCES "chapter_comments" ("id");

ALTER TABLE "chapter_comment_likes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_bookmarks" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_bookmarks" ADD FOREIGN KEY ("publication_id") REFERENCES "publications" ("id");

ALTER TABLE "user_bookmarks" ADD FOREIGN KEY ("last_chapter_interacted_id") REFERENCES "chapters" ("id");

ALTER TABLE "user_bookmarks" ADD FOREIGN KEY ("status_id") REFERENCES "user_bookmark_status" ("id");

ALTER TABLE "chapter_reviews" ADD FOREIGN KEY ("chapter_id") REFERENCES "chapters" ("id");

ALTER TABLE "chapter_reviews" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "chapter_review_likes" ADD FOREIGN KEY ("review_id") REFERENCES "chapter_reviews" ("id");

ALTER TABLE "chapter_review_likes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

-- Create two computed columns for the unique constraint
ALTER TABLE "suggestions" ADD COLUMN "smallest_publication_id" bigint GENERATED ALWAYS AS (LEAST("publication_1_id", "publication_2_id")) STORED,
ADD COLUMN "largest_publication_id" bigint GENERATED ALWAYS AS (GREATEST("publication_1_id", "publication_2_id")) STORED;

-- Add the unique constraint on these computed columns
ALTER TABLE "suggestions" ADD CONSTRAINT "unique_publication_pair"
UNIQUE ("smallest_publication_id", "largest_publication_id");

-- Add unique constraint for user_bookmarks
ALTER TABLE "user_bookmarks" ADD CONSTRAINT "user_bookmarks_user_id_publication_id_key" 
UNIQUE ("user_id", "publication_id");
