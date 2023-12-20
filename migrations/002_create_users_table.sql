-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    username character varying(255) NOT NULL,
    email character varying(255) NOT NULL UNIQUE
);

-- +migrate Down
DROP TABLE IF EXISTS users;