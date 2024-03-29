DROP TABLE IF EXISTS chats CASCADE;
CREATE TABLE chats
(
    ID        BIGSERIAL NOT NULL PRIMARY KEY,
    supportID BIGINT    NULL,
    userID    BIGINT    NULL,
    FOREIGN KEY (supportID) REFERENCES supports (ID) ON DELETE CASCADE
);