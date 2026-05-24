-- +goose Up
CREATE TABLE IF NOT EXISTS achievement_groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    position INT DEFAULT 0,
    is_hidden BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS achievements (
    id SERIAL PRIMARY KEY,
    category_id INT REFERENCES achievement_groups(id) ON DELETE CASCADE,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    unlock_method TEXT NOT NULL, -- Study more if this is really the best way to do it. 'deal_damage', 'kill_enemies', 'end_quest'
    unlock_requirement_amount INT NULL, -- Damage Value, Enemies Amount
    required_achievement_id  INT NULL,
    required_quest_id INT NULL,
    is_secret BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS user_achievements (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    achievement_id INT REFERENCES achievements(id) ON DELETE CASCADE,
    current_progress INT DEFAULT 0,
    is_complete BOOLEAN DEFAULT false,
    unlocked_at TIMESTAMP DEFAULT NOW(),
    last_updated TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, achievement_id)
);

CREATE INDEX IF NOT EXISTS idx_achievements_method ON achievements(unlock_method);
CREATE INDEX IF NOT EXISTS idx_achievement_category ON achievements(category_id);
CREATE INDEX IF NOT EXISTS idx_user_achievement_complete ON user_achievements(user_id, is_complete);
CREATE INDEX IF NOT EXISTS idx_user_achievement_progress ON user_achievements(current_progress);

INSERT INTO achievement_groups (name, position, is_hidden)
VALUES('CODENAME // DEVELOPMENT', 0, false);

INSERT INTO achievement_groups (name, position, is_hidden)
VALUES('Beyond the Digital Ocean', 1, false);

INSERT INTO achievement_groups (name, position, is_hidden)
VALUES('Digital Oddities', 2, true);

INSERT INTO achievements (category_id, name, description, unlock_method, unlock_requirement_amount)
VALUES (1, 'DEBUG MISSION', 'LOREM IPSUN IDK LOL', 'cheating', 99999);

INSERT INTO achievements (category_id, name, description, unlock_method, unlock_requirement_amount, required_achievement_id)
VALUES (1, 'YET ANOTHER DEBG MISSION', 'MEWGENICS IS PEAK', 'end_quest', 99999, 1);

INSERT INTO achievements (category_id, name, description, unlock_method, required_quest_id)
VALUES (2, 'Memories from the past', 'You woke up into this world.', 'natural', 1);


-- +goose Down
DROP TABLE IF EXISTS user_achievements;
DROP TABLE IF EXISTS achievements;
DROP TABLE IF EXISTS achievement_groups;

DROP INDEX IF EXISTS idx_user_achievement_progress;
DROP INDEX IF EXISTS idx_user_achievement_complete;
DROP INDEX IF EXISTS idx_achievement_category;
DROP INDEX IF EXISTS idx_achievements_method;
