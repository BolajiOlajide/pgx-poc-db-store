-- +migrate Up
CREATE TABLE IF NOT EXISTS people (
    id SERIAL PRIMARY KEY,
    user_id uuid REFERENCES users(id) UNIQUE NOT NULL
);


-- +migrate Down
DROP TABLE IF EXISTS people;