-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_digimon (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    base_id INTEGER REFERENCES base_digimon(id),

    -- General Info
    nickname TEXT,
    level INTEGER DEFAULT 1,
    exp INTEGER DEFAULT 0,
    friendship INTEGER DEFAULT 0,

    -- Flags
    is_starter BOOLEAN DEFAULT false,
    is_locked BOOLEAN DEFAULT false,

    -- Stats
    current_health INTEGER,
    current_mana INTEGER,
    health INTEGER,
    mana INTEGER,
    attack INTEGER,
    defense INTEGER,
    speed INTEGER,

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_team (
    user_id UUID REFERENCES users(id),
    digimon_id UUID REFERENCES user_digimon(id),
    team_slot INTEGER CHECK (team_slot BETWEEN 1 AND 3),
    PRIMARY KEY (user_id, team_slot),
    UNIQUE(user_id, digimon_id)
);

CREATE TABLE IF NOT EXISTS user_digibank (
    user_id UUID REFERENCES users(id),
    digimon_id UUID REFERENCES user_digimon(id),
    bank_slot INTEGER CHECK (bank_slot BETWEEN 1 AND 3),
    PRIMARY KEY (user_id, bank_slot)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_team;
DROP TABLE IF EXISTS user_digibank;
DROP TABLE IF EXISTS user_digimon;
-- +goose StatementEnd
