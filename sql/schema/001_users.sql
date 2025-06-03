-- +goose Up
CREATE TABLE users(
    id uuid PRIMARY KEY NOT NULL ,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashed_password text not null default 'unset'
);

-- +goose Down
DROP TABLE users;