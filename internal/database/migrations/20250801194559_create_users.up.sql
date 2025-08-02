CREATE TABLE users
(
    id         CHAR(36) PRIMARY KEY,
    username   TEXT      NOT NULL UNIQUE,
    email      TEXT      NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL
);