-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chapters (
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


CREATE TABLE IF NOT EXISTS episodes (
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

CREATE TABLE IF NOT EXISTS stages (
    id SERIAL PRIMARY KEY,
    chapter_id INTEGER REFERENCES chapters(id),
    episode_id INTEGER REFERENCES episodes(id),
    stage_number INTEGER NOT NULL,
    stage_name TEXT NOT NULL,
    description TEXT,

    -- Requierements
    defeat_stage_id INTEGER REFERENCES stages(id) DEFAULT NULL,
    tamer_level INTEGER DEFAULT 0,

    -- Awards
    drop_bits INTEGER DEFAULT 0,
    drop_exp INTEGER DEFAULT 0,
    drop_sexp INTEGER DEFAULT 0,

    -- Config
    max_waves INTEGER DEFAULT 1 check (max_waves BETWEEN 1 AND 3),
    stage_type TEXT DEFAULT 'battle'
        CHECK (stage_type IN ('battle', 'boss_battle', 'dialogue', 'cutscene')),
    
    -- Flags System
    trigger_flag_on_dialogue_end TEXT DEFAULT NULL,

    UNIQUE(episode_id, stage_number)
);

INSERT INTO stages (episode_id, chapter_id, stage_number, stage_name, description, defeat_stage_id, stage_type)
VALUES (1, 1, 1, 'Testando o teste', 'Teste', NULL, 'dialogue');
INSERT INTO stages (episode_id, chapter_id, stage_number, stage_name, description, defeat_stage_id)
VALUES (1, 1, 2, 'Testando o teste testado', 'Teste', 1);
INSERT INTO stages (episode_id, chapter_id, stage_number, stage_name, description, defeat_stage_id)
VALUES (1, 1, 3, 'funcionou?', 'Teste', 2);


-- Rewrite later
CREATE TABLE IF NOT EXISTS stage_waves (
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


CREATE TABLE IF NOT EXISTS wave_enemies (
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wave_enemies;
DROP TABLE IF EXISTS stage_waves;
DROP TABLE IF EXISTS stages;
DROP TABLE IF EXISTS episodes;
DROP TABLE IF EXISTS chapters;
-- +goose StatementEnd

