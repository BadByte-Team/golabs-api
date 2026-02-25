CREATE TABLE challenges (
    id          CHAR(36) PRIMARY KEY,
    event_id    CHAR(36) NOT NULL,
    title       VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category    VARCHAR(64) NOT NULL,
    points      INT NOT NULL DEFAULT 0,
    difficulty  VARCHAR(32) NOT NULL DEFAULT 'medium',
    visible     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  DATETIME NOT NULL,
    updated_at  DATETIME NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id)
);

CREATE TABLE flags (
    id           CHAR(36) PRIMARY KEY,
    challenge_id CHAR(36) NOT NULL UNIQUE,
    hash         CHAR(64) NOT NULL,       -- SHA-256 hex, plain-text NUNCA persiste
    created_at   DATETIME NOT NULL,
    FOREIGN KEY (challenge_id) REFERENCES challenges(id)
);

CREATE TABLE solves (
    id            CHAR(36) PRIMARY KEY,
    challenge_id  CHAR(36) NOT NULL,
    event_team_id CHAR(36) NOT NULL,
    user_id       CHAR(36) NOT NULL,
    solved_at     DATETIME NOT NULL,
    UNIQUE KEY uq_team_challenge (challenge_id, event_team_id),
    FOREIGN KEY (challenge_id) REFERENCES challenges(id),
    FOREIGN KEY (event_team_id) REFERENCES event_teams(id)
);
