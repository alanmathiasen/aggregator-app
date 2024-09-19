-- Remove foreign key constraints
ALTER TABLE "publication_types" DROP CONSTRAINT IF EXISTS "publication_types_id_fkey";
ALTER TABLE "publications" DROP CONSTRAINT IF EXISTS "publications_id_fkey";
ALTER TABLE "publication_reviews" DROP CONSTRAINT IF EXISTS "publication_reviews_id_fkey";
ALTER TABLE "users" DROP CONSTRAINT IF EXISTS "users_id_fkey";
ALTER TABLE "chapters" DROP CONSTRAINT IF EXISTS "chapters_publication_id_fkey";
ALTER TABLE "chapter_likes" DROP CONSTRAINT IF EXISTS "chapter_likes_chapter_id_fkey";
ALTER TABLE "chapter_comments" DROP CONSTRAINT IF EXISTS "chapter_comments_chapter_id_fkey";
ALTER TABLE "chapter_comment_likes" DROP CONSTRAINT IF EXISTS "chapter_comment_likes_comment_id_fkey";
ALTER TABLE "user_bookmarks" DROP CONSTRAINT IF EXISTS "user_bookmarks_status_id_fkey";
ALTER TABLE "chapter_reviews" DROP CONSTRAINT IF EXISTS "chapter_reviews_id_fkey";
ALTER TABLE "chapter_review_likes" DROP CONSTRAINT IF EXISTS "chapter_review_likes_is_like_fkey";

-- Remove unique constraints
ALTER TABLE "user_bookmarks" DROP CONSTRAINT IF EXISTS "user_bookmarks_user_id_publication_id_key";
ALTER TABLE "suggestions" DROP CONSTRAINT IF EXISTS "unique_publication_pair";

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
DROP TABLE IF EXISTS "publication_types";
DROP TABLE IF EXISTS "publications";
DROP TABLE IF EXISTS "users";