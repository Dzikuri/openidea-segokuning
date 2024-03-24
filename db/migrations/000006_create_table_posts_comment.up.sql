-- Create Table
CREATE TABLE IF NOT EXISTS post_comments (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "post_id" uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    "user_id" uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    "comment" text NOT NULL,
    "created_at" timestamptz(6),
    "updated_at" timestamptz(6)
);

CREATE INDEX IF NOT EXISTS "idx_post_comment_posts_id" ON "public"."post_comments"("post_id");

CREATE INDEX IF NOT EXISTS "idx_post_comment_users_id" ON "public"."post_comments"("user_id");

CREATE INDEX IF NOT EXISTS "idx_post_comment_posts_users_id" ON "public"."post_comments"("post_id", "user_id");