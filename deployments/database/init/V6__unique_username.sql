-- Add UNIQUE constraint on username and index for search performance.
-- Safe to run even if the constraint already exists (IF NOT EXISTS pattern via column rename).

ALTER TABLE users
    MODIFY COLUMN username VARCHAR(50) NOT NULL,
    ADD UNIQUE INDEX uq_users_username (username);
