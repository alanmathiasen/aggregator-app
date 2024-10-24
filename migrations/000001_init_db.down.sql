-- First, drop all foreign key constraints
ALTER TABLE IF EXISTS "publications" DROP CONSTRAINT IF EXISTS "publications_type_id_fkey";
ALTER TABLE IF EXISTS "publication_reviews" DROP CONSTRAINT IF EXISTS "publication_reviews_publication_id_fkey";
ALTER TABLE IF EXISTS "publication_reviews" DROP CONSTRAINT IF EXISTS "publication_reviews_user_id_fkey";
ALTER TABLE IF EXISTS "chapters" DROP CONSTRAINT IF EXISTS "chapters_publication_id_fkey";
ALTER TABLE IF EXISTS "chapter_likes" DROP CONSTRAINT IF EXISTS "chapter_likes_chapter_id_fkey";
ALTER TABLE IF EXISTS "chapter_likes" DROP CONSTRAINT IF EXISTS "chapter_likes_user_id_fkey";
ALTER TABLE IF EXISTS "chapter_comments" DROP CONSTRAINT IF EXISTS "chapter_comments_chapter_id_fkey";
ALTER TABLE IF EXISTS "chapter_comments" DROP CONSTRAINT IF EXISTS "chapter_comments_user_id_fkey";
ALTER TABLE IF EXISTS "chapter_comment_likes" DROP CONSTRAINT IF EXISTS "chapter_comment_likes_comment_id_fkey";
ALTER TABLE IF EXISTS "chapter_comment_likes" DROP CONSTRAINT IF EXISTS "chapter_comment_likes_user_id_fkey";
ALTER TABLE IF EXISTS "user_bookmarks" DROP CONSTRAINT IF EXISTS "user_bookmarks_user_id_fkey";
ALTER TABLE IF EXISTS "user_bookmarks" DROP CONSTRAINT IF EXISTS "user_bookmarks_publication_id_fkey";
ALTER TABLE IF EXISTS "user_bookmarks" DROP CONSTRAINT IF EXISTS "user_bookmarks_status_id_fkey";
ALTER TABLE IF EXISTS "chapter_reviews" DROP CONSTRAINT IF EXISTS "chapter_reviews_chapter_id_fkey";
ALTER TABLE IF EXISTS "chapter_reviews" DROP CONSTRAINT IF EXISTS "chapter_reviews_user_id_fkey";
ALTER TABLE IF EXISTS "chapter_review_likes" DROP CONSTRAINT IF EXISTS "chapter_review_likes_review_id_fkey";
ALTER TABLE IF EXISTS "chapter_review_likes" DROP CONSTRAINT IF EXISTS "chapter_review_likes_user_id_fkey";
ALTER TABLE IF EXISTS "review_likes" DROP CONSTRAINT IF EXISTS "review_likes_review_id_fkey";
ALTER TABLE IF EXISTS "review_likes" DROP CONSTRAINT IF EXISTS "review_likes_user_id_fkey";
ALTER TABLE IF EXISTS "suggestions" DROP CONSTRAINT IF EXISTS "suggestions_publication_id_fkey";
ALTER TABLE IF EXISTS "suggestions" DROP CONSTRAINT IF EXISTS "suggestions_suggested_publication_id_fkey";
ALTER TABLE IF EXISTS "publication_genres" DROP CONSTRAINT IF EXISTS "publication_genres_publication_id_fkey";
ALTER TABLE IF EXISTS "publication_genres" DROP CONSTRAINT IF EXISTS "publication_genres_genre_id_fkey";

-- Remove unique constraints
ALTER TABLE IF EXISTS "user_bookmarks" DROP CONSTRAINT IF EXISTS "user_bookmarks_user_id_publication_id_key";
ALTER TABLE IF EXISTS "suggestions" DROP CONSTRAINT IF EXISTS "unique_publication_pair";

-- Drop tables in reverse order of creation
DROP TABLE IF EXISTS "chapter_review_likes";
DROP TABLE IF EXISTS "chapter_reviews";
DROP TABLE IF EXISTS "suggestions";
DROP TABLE IF EXISTS "user_bookmark_status";
DROP TABLE IF EXISTS "user_bookmarks";
DROP TABLE IF EXISTS "chapter_comment_likes";
DROP TABLE IF EXISTS "chapter_comments";
DROP TABLE IF EXISTS "chapter_likes";
DROP TABLE IF EXISTS "chapters";
DROP TABLE IF EXISTS "review_likes";
DROP TABLE IF EXISTS "publication_reviews";
DROP TABLE IF EXISTS "publication_genres";
DROP TABLE IF EXISTS "publications";
DROP TABLE IF EXISTS "publication_types";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "genres";

-- Drop the vector extension if it was added
DROP EXTENSION IF EXISTS vector;

-- Finally, drop the image_url column from publications if it still exists
ALTER TABLE IF EXISTS publications DROP COLUMN IF EXISTS image_url;
