DROP INDEX IF EXISTS index_users_on_total_friend;

ALTER TABLE
    users DROP COLUMN IF EXISTS total_friends;