-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_completed_stages (
    user_id UUID REFERENCES users(id),
    stage_id INTEGER REFERENCES stages(id),
    completed_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY(user_id, stage_id)
);

CREATE TABLE IF NOT EXISTS user_chapter_progress (
    user_id UUID REFERENCES users(id),
    chapter_id INTEGER REFERENCES chapters(id),
    is_unlocked BOOLEAN DEFAULT false,
    is_beaten BOOLEAN DEFAULT false,
    completed_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY(user_id, chapter_id)   
);

CREATE TABLE IF NOT EXISTS user_episode_progress (
    user_id UUID REFERENCES users(id),
    episode_id INTEGER REFERENCES episodes(id),
    is_unlocked BOOLEAN DEFAULT false,
    is_beaten BOOLEAN DEFAULT false,
    completed_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY(user_id, episode_id)   
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_completed_stages;
DROP TABLE IF EXISTS user_chapter_progress;
DROP TABLE IF EXISTS user_episode_progress;
-- +goose StatementEnd
