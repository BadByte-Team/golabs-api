-- V7: Add performance indexes for high-traffic query patterns
-- These indexes prevent full table scans on the most frequent read paths.

-- Challenges: filter by event + visibility (list endpoint) and by event + category (filter endpoint)
CREATE INDEX IF NOT EXISTS idx_challenges_event_visible
    ON challenges (event_id, visible);

CREATE INDEX IF NOT EXISTS idx_challenges_event_category
    ON challenges (event_id, category);

-- Solves: first-blood and solve-count queries
CREATE INDEX IF NOT EXISTS idx_solves_challenge_solved_at
    ON solves (challenge_id, solved_at);

-- Event teams: leaderboard query (list by event)
CREATE INDEX IF NOT EXISTS idx_event_teams_event_score
    ON event_teams (event_id, score DESC);

-- Event team members: exists check and member list
CREATE INDEX IF NOT EXISTS idx_etm_team_id
    ON event_team_members (event_team_id);

CREATE INDEX IF NOT EXISTS idx_etm_event_user
    ON event_team_members (user_id);
