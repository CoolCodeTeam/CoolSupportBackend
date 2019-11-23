DROP TABLE IF EXISTS supports CASCADE;
CREATE TABLE supports
(
    ID           BIGSERIAL    NOT NULL PRIMARY KEY,
    email        VARCHAR(128) NOT NULL UNIQUE,
    password BYTEA        NOT NULL
);