CREATE TABLE "publications" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "description" text NOT NULL,
  "image" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "chapters" (
  "id" bigserial PRIMARY KEY,
  "publication_id" bigint NOT NULL,
  "title" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "publication_sources" (
  "id" bigserial PRIMARY KEY,
  "publication_id" bigint,
  "link" varchar NOT NULL,
  "source" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "avatar" varchar,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz DEFAULT (now())
);

CREATE TABLE "user_publication_follows" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "publication_id" bigint NOT NULL,
  "chapter_id" bigint,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz DEFAULT (now())
);

CREATE TABLE "user_publication_ratings" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "publication_id" bigint NOT NULL,
  "rating" float NOT NULL,
  "text" text,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz DEFAULT (now())
);

CREATE TABLE "user_chapter_ratings" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint,
  "chapter_id" bigint,
  "rating" float,
  "text" text,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz DEFAULT (now())
);

CREATE TABLE "comments" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "publication_id" bigint,
  "chapter_id" bigint,
  "parent_comment_id" bigint,
  "content" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "chapters" ("publication_id");

CREATE INDEX ON "publication_sources" ("publication_id");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "user_publication_ratings" ("publication_id");

CREATE INDEX ON "user_chapter_ratings" ("chapter_id");

CREATE INDEX ON "comments" ("publication_id");

CREATE INDEX ON "comments" ("chapter_id");

ALTER TABLE "chapters" ADD FOREIGN KEY ("publication_id") REFERENCES "publications" ("id");

ALTER TABLE "publication_sources" ADD FOREIGN KEY ("publication_id") REFERENCES "publications" ("id");

ALTER TABLE "user_publication_ratings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_publication_ratings" ADD FOREIGN KEY ("publication_id") REFERENCES "publications" ("id");

ALTER TABLE "user_chapter_ratings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_chapter_ratings" ADD FOREIGN KEY ("chapter_id") REFERENCES "chapters" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("publication_id") REFERENCES "publications" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("chapter_id") REFERENCES "chapters" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("parent_comment_id") REFERENCES "comments" ("id");
