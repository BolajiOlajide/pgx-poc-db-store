-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE
);

-- +migrate Down
DROP TABLE IF EXISTS users;