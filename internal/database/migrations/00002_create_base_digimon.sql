-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS base_digimon (
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

INSERT INTO base_digimon(species, base_health, base_attack, base_defense, base_mana, base_speed, form, attribute, family, element) 
VALUES ('Chronomon: Holy Mode', 999999, 999999, 999999, 999999, 999999, 'Ultimate', 'Vaccine', 'Virus Busters', 'Light');

INSERT INTO base_digimon(species, base_health, base_attack, base_defense, base_mana, base_speed, form, attribute, family, element) 
VALUES ('Monodramon', 150, 120, 130, 100, 100, 'Rookie', 'Vaccine', 'Dragon''s Roar', 'Fire');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS base_digimon;
-- +goose StatementEnd
