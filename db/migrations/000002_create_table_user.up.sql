-- NOTE Create table users
CREATE TABLE IF NOT EXISTS "public"."users" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    "email" varchar(255) UNIQUE,
    "phone" varchar(15) UNIQUE,
    "name" varchar(50) NOT NULL,
    "password" varchar(255) NOT NULL,
    "image_url" text DEFAULT '',
    "created_at" timestamptz(6),
    "updated_at" timestamptz(6)
);

-- NOTE Index columns
CREATE INDEX idx_users_name ON "public"."users" ("name");

CREATE INDEX idx_users_email ON "public"."users" ("email");

CREATE INDEX idx_users_phone ON "public"."users" ("phone");

-- NOTE Composite index columns
CREATE INDEX idx_users_name_email ON "public"."users" ("name", "email");