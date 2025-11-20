-- name: GetStageById :one
SELECT s.*, e.episode_name, c.chapter_name
FROM stages s
JOIN episodes e ON s.episode_id = e.id
JOIN chapters c ON s.chapter_id = c.id
WHERE s.id = $1;

-- name: GetStageWaves :many
SELECT sw.*, we.*, bd.*
FROM stage_waves sw
LEFT JOIN wave_enemies we ON sw.id = we.id
LEFT JOIN base_digimon bd ON we.enemy_id = bd.id
WHERE sw.id = $1
ORDER BY sw.wave_number, we.enemy_slot;

-- name: GetAvailableStages :many
SELECT s.*, e.episode_name, c.chapter_name
FROM stages s
JOIN episodes e ON s.episode_id = e.id
JOIN chapters c ON e.chapter_id = c.id
WHERE s.defeat_stage_id IS NULL
    OR EXISTS (
        SELECT 1 FROM user_completed_stages 
        WHERE user_id = $1 AND stage_id = s.defeat_stage_id
    )
ORDER BY c.chapter_number, e.episode_number, s.stage_number;

-- name: GetAvailableEpisodesByChapterId :many
SELECT e.*
FROM episodes e
WHERE e.id = 1 AND e.chapter_id = $1
    OR EXISTS (
        SELECT 1 FROM user_episode_progress uep
        JOIN episodes prev_ep ON prev_ep.episode_number = e.episode_number - 1 
        AND prev_ep.chapter_id = e.chapter_id
        WHERE uep.user_id = $2 AND uep.episode_id = prev_ep.id
    )
ORDER BY e.episode_number;

-- name: GetAvailableStagesByEpisodeId :many
SELECT s.*
FROM stages s
WHERE s.episode_id = $1
  AND (
    s.defeat_stage_id IS NULL 
    OR EXISTS (
      SELECT 1 FROM user_completed_stages 
      WHERE user_id = $2
        AND stage_id = s.defeat_stage_id
    )
  )
ORDER BY s.stage_number;