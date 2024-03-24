-- Create Table
CREATE TABLE IF NOT EXISTS posts (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "user_id" uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    "content" text NOT NULL,
    "tags" text [] NOT NULL,
    "created_at" timestamptz(6),
    "updated_at" timestamptz(6)
);

CREATE INDEX IF NOT EXISTS "idx_posts_user_id" ON "public"."posts"("user_id");

CREATE INDEX IF NOT EXISTS "idx_posts_tags" ON "public"."posts"("tags");