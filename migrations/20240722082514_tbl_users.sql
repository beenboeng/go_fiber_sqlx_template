-- +goose Up
CREATE TABLE IF NOT EXISTS tbl_users (
    "id" SERIAL PRIMARY KEY,
    "first_name" TEXT,
    "last_name" TEXT,
    "username" TEXT NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    "email" VARCHAR(100) NOT NULL,
    "phone" TEXT,
    "login_session" TEXT,
    "created_by" int,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_by" int,
    "updated_at" TIMESTAMP(3),
    "deleted_by" int,
    "deleted_at" TIMESTAMP(3)
);

-- +goose Down
DROP TABLE IF EXISTS tbl_users;

