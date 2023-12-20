-- +migrate Up
CREATE TABLE IF NOT EXISTS people (id int);


-- +migrate Down
DROP TABLE IF EXISTS people;