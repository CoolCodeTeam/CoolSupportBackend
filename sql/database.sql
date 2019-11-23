DROP DATABASE IF EXISTS supportdatabase;
CREATE DATABASE supportdatabase;


DROP TABLE IF EXISTS supports CASCADE;
CREATE TABLE supports
(
    ID           BIGSERIAL    NOT NULL PRIMARY KEY,
    email        VARCHAR(128) NOT NULL UNIQUE,
    password BYTEA        NOT NULL
);


DROP TABLE IF EXISTS chats CASCADE;
CREATE TABLE chats
(
    ID        BIGSERIAL NOT NULL PRIMARY KEY,
    supportID BIGINT    NULL,
    userID    BIGINT    NULL,
    FOREIGN KEY (supportID) REFERENCES supports (ID) ON DELETE CASCADE
);


DROP TABLE IF EXISTS messages CASCADE;
CREATE TABLE messages
(
    ID          BIGSERIAL NOT NULL PRIMARY KEY,
    body        TEXT      NOT NULL,
    chatID      BIGINT    NOT NULL,
    messageTime TIMESTAMP,
    isSupport   BOOLEAN,
    FOREIGN KEY (chatID) REFERENCES chats (ID)
)