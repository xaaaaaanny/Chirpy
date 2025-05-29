-- +goose Up
CREATE TABLE chirps(
    id uuid PRIMARY KEY NOT NULL ,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    body TEXT NOT NULL,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE NOT NULL                  
);

-- +goose Down
DROP TABLE chirps;