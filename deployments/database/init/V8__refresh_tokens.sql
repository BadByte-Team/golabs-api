CREATE TABLE IF NOT EXISTS refresh_tokens (
    id         CHAR(36)     NOT NULL PRIMARY KEY,
    user_id    CHAR(36)     NOT NULL,
    token_hash CHAR(64)     NOT NULL UNIQUE, -- SHA-256 hex of the raw token
    expires_at TIMESTAMP    NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP    NULL     DEFAULT NULL,
    CONSTRAINT fk_rt_user FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_rt_token_hash ON refresh_tokens (token_hash);
CREATE INDEX IF NOT EXISTS idx_rt_user_id    ON refresh_tokens (user_id);
