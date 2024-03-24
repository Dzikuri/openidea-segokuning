-- NOTE Create table users
CREATE TABLE IF NOT EXISTS "public"."friends" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "user_id" uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    "follow_user_id" uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    "created_at" timestamptz(6),
    "updated_at" timestamptz(6)
);

-- NOTE Index columns
CREATE INDEX IF NOT EXISTS idx_friends_user_id ON "public"."friends" ("user_id");

CREATE INDEX IF NOT EXISTS idx_friends_follow_user_id ON "public"."friends" ("follow_user_id");