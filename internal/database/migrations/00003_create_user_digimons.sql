-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_digimon (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    base_id INTEGER REFERENCES base_digimon(id),

    -- General Info
    species TEXT,
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

CREATE OR REPLACE FUNCTION set_initial_digimon_stats()
RETURNS TRIGGER AS $$
BEGIN
    SELECT 
        species, base_health, base_mana, base_attack, base_defense, base_speed
    INTO
        NEW.species, NEW.health, NEW.mana, NEW.attack, NEW.defense, NEW.speed
    FROM base_digimon
    WHERE id = NEW.base_id;

    NEW.current_health := NEW.health;
    NEW.current_mana := NEW.mana;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_initial_digimon_stats
    BEFORE INSERT ON user_digimon
    FOR EACH ROW 
    EXECUTE FUNCTION set_initial_digimon_stats();

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
DROP FUNCTION IF EXISTS set_initial_digimon_stats;
DROP TRIGGER IF EXISTS trigger_set_initial_digimon_stats ON user_digimon;
-- +goose StatementEnd
