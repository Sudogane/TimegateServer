-- name: GetBaseDigimon :one
SELECT * FROM base_digimon WHERE id = $1;

-- name: GetBaseDigimonBySpecies :one
SELECT * FROM base_digimon WHERE species = $1;

-- name: GetAllBaseDigimon :many
SELECT * FROM base_digimon ORDER BY species;

-- name: CreateUserDigimon :one
INSERT INTO user_digimon (user_id, base_id, is_starter, is_locked) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetUserDigimon :one
SELECT * FROM user_digimon WHERE id = $1 AND user_id = $2;

-- name: GetAllUserDigimon :many
SELECT * FROM user_digimon WHERE id = $1;

-- name: GetUserDigimonByStarterFlag :one
SELECT * FROM user_digimon WHERE user_id = $1 AND is_starter = true LIMIT 1;