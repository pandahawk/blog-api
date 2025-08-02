CREATE TABLE posts
(
    id         CHAR(36) PRIMARY KEY,
    title      TEXT      NOT NULL,
    content    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id    CHAR(36)  NOT NULL,
    CONSTRAINT fk_users_posts
        FOREIGN KEY (user_id)
            REFERENCES USERS (id)
            ON DELETE CASCADE
);