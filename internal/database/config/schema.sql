CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    crated_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    online BOOLEAN DEFAULT true
);

CREATE TABLE user_resources (
    user_id UUID PRIMARY KEY REFERENCES users(id),

    -- General Info
    level INTEGER DEFAULT 1,
    exp INTEGER DEFAULT 0,
    stamina_current INTEGER DEFAULT 100,
    stamina_max INTEGER DEFAULT 100,
    
    -- Currencies
    bits BIGINT DEFAULT 1500,
    yen BIGINT DEFAULT 0
);

CREATE TABLE user_inventory (
    id SERIAL,
    user_id UUID PRIMARY KEY REFERENCES users(id),
    item_name TEXT NOT NULL,
    quantity INTEGER DEFAULT 1
);

CREATE TABLE user_digimon (
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

CREATE TABLE user_team (
    user_id UUID REFERENCES users(id),
    digimon_id UUID REFERENCES user_digimon(id),
    team_slot INTEGER CHECK (team_slot BETWEEN 1 AND 3),
    PRIMARY KEY (user_id, team_slot),
    UNIQUE(user_id, digimon_id)
);

CREATE TABLE user_digibank (
    user_id UUID REFERENCES users(id),
    digimon_id UUID REFERENCES user_digimon(id),
    bank_slot INTEGER CHECK (bank_slot BETWEEN 1 AND 3),
    PRIMARY KEY (user_id, bank_slot)
);

CREATE TABLE base_digimon (
    id SERIAL PRIMARY KEY,
    species TEXT UNIQUE NOT NULL,
    base_health INTEGER NOT NULL,
    base_attack INTEGER NOT NULL,
    base_defense INTEGER NOT NULL,
    base_mana INTEGER NOT NULL,
    base_speed INTEGER NOT NULL,
    form TEXT NOT NULL,
    attribute TEXT NOT NULL,
    family TEXT NOT NULL,
    element TEXT NOT NULL
);

INSERT INTO base_digimon(species, base_health, base_attack, base_defense, base_mana, base_speed, form, attribute, family, element) 
VALUES ('Morphomon', 100, 100, 100, 100, 100, 'Rookie', 'Vaccine', 'Nature Spirits', 'Light');

CREATE TABLE chapters (
    id SERIAL PRIMARY KEY,
    chapter_name TEXT NOT NULL,
    chapter_number TEXT NOT NULL UNIQUE,
    description TEXT,
    is_unlocked BOOLEAN DEFAULT false
);

INSERT INTO chapters (chapter_name, description, is_unlocked, chapter_number)
VALUES ('Alpha', 'Uh, testing phase i guess?', true, 1);
INSERT INTO chapters (chapter_name, description, is_unlocked, chapter_number)
VALUES ('Omega', 'it worked i guess?', true, 2);

CREATE TABLE episodes (
    id SERIAL PRIMARY KEY,
    chapter_id INTEGER REFERENCES chapters(id),
    episode_number INTEGER NOT NULL,
    episode_name TEXT NOT NULL,

    -- Adicionar depois a variavel das cenas de dialogo
    UNIQUE(chapter_id, episode_number)
);

INSERT INTO episodes (chapter_id, episode_number, episode_name)
VALUES (1, 1, 'Testing');
INSERT INTO episodes (chapter_id, episode_number, episode_name)
VALUES (1, 2, 'Phase');

CREATE TABLE stages (
    id SERIAL PRIMARY KEY,
    chapter_id INTEGER REFERENCES chapters(id),
    episode_id INTEGER REFERENCES episodes(id),
    stage_number INTEGER NOT NULL,
    stage_name TEXT NOT NULL,
    description TEXT,

    -- Requierements
    defeat_stage_id INTEGER REFERENCES stages(id),
    tamer_level INTEGER DEFAULT 0,

    -- Awards
    drop_bits INTEGER DEFAULT 0,
    drop_exp INTEGER DEFAULT 0,
    drop_sexp INTEGER DEFAULT 0,

    -- Config
    max_waves INTEGER DEFAULT 1 check (max_waves BETWEEN 1 AND 3),

    UNIQUE(episode_id, stage_number)
);

INSERT INTO stages (episode_id, chapter_id, stage_number, stage_name, description, defeat_stage_id)
VALUES (1, 1, 1, 'Testando o teste', 'Teste');
INSERT INTO stages (episode_id, stage_number, stage_name, description, defeat_stage_id)
VALUES (1, 1, 2, 'Testando o teste testado', 'Teste', 1);
INSERT INTO stages (episode_id, stage_number, stage_name, description, defeat_stage_id)
VALUES (1, 1, 3, 'funcionou?', 'Teste', 2);

CREATE TABLE stage_waves (
    id SERIAL PRIMARY KEY,
    stage_id INTEGER REFERENCES stages(id),
    wave_number INTEGER NOT NULL CHECK (wave_number BETWEEN 1 AND 3),
    UNIQUE(stage_id, wave_number)
);

INSERT INTO stage_waves (stage_id, wave_number)
VALUES (1, 1);
INSERT INTO stage_waves (stage_id, wave_number)
VALUES (2, 1);
INSERT INTO stage_waves (stage_id, wave_number)
VALUES (3, 1);

CREATE TABLE wave_enemies (
    id SERIAL PRIMARY KEY,
    wave_id INTEGER REFERENCES stage_waves(id),
    
    enemy_slot INTEGER NOT NULL CHECK (enemy_slot BETWEEN 1 AND 3),
    enemy_id INTEGER REFERENCES base_digimon(id),
    enemy_level INTEGER DEFAULT 1,

    UNIQUE(wave_id, enemy_slot)
);

INSERT INTO wave_enemies (wave_id, enemy_slot, enemy_id, enemy_level)
VALUES (1, 1, 1, 1);
INSERT INTO wave_enemies (wave_id, enemy_slot, enemy_id, enemy_level)
VALUES (2, 2, 1, 2);
INSERT INTO wave_enemies (wave_id, enemy_slot, enemy_id, enemy_level)
VALUES (2, 1, 1, 1);
INSERT INTO wave_enemies (wave_id, enemy_slot, enemy_id, enemy_level)
VALUES (3, 1, 1, 10);
INSERT INTO wave_enemies (wave_id, enemy_slot, enemy_id, enemy_level)
VALUES (3, 2, 1, 20);

-- Future
CREATE TABLE stage_item_drops (
    id SERIAL PRIMARY KEY,
    stage_id INTEGER REFERENCES stages(id)
);

CREATE TABLE user_completed_stages (
    user_id UUID REFERENCES users(id),
    stage_id INTEGER REFERENCES stages(id),
    completed_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY(user_id, stage_id)
);

CREATE TABLE user_chapter_progress (
    user_id UUID REFERENCES users(id),
    chapter_id INTEGER REFERENCES chapters(id),
    is_unlocked BOOLEAN DEFAULT false,
    is_beaten BOOLEAN DEFAULT false,
    completed_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY(user_id, chapter_id)   
);

CREATE TABLE user_episode_progress (
    user_id UUID REFERENCES users(id),
    episode_id INTEGER REFERENCES episodes(id),
    is_unlocked BOOLEAN DEFAULT false,
    is_beaten BOOLEAN DEFAULT false,
    completed_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY(user_id, episode_id)   
);

CREATE INDEX CONCURRENTLY idx_users_id ON users(id);
CREATE INDEX CONCURRENTLY idx_users_username ON users(username);
CREATE INDEX CONCURRENTLY idx_user_resources_id ON user_resources(id);
CREATE INDEX CONCURRENTLY idx_user_team_user_id ON user_team(user_id);
CREATE INDEX CONCURRENTLY idx_user_digimons_user_id ON user_digimon(user_id);
CREATE INDEX CONCURRENTLY idx_user_digibank_user_id ON user_digibank(user_id);
