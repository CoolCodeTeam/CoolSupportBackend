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