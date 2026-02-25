CREATE TABLE event_teams (
    id CHAR(36) PRIMARY KEY,
    event_id CHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    join_secret_hash VARCHAR(255) NOT NULL,
    score INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(event_id, name),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

CREATE TABLE event_team_members (
    event_team_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    role ENUM('owner','member') NOT NULL,
    joined_at TIMESTAMP NOT NULL,
    PRIMARY KEY (event_team_id, user_id),
    FOREIGN KEY (event_team_id) REFERENCES event_teams(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);