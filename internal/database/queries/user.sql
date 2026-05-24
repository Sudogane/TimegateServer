-- name: CreateUser :one
INSERT INTO users (
    username, password_hash
) VALUES ($1, $2) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: CheckIfUsernameIsTaken :one
SELECT EXISTS (SELECT 1 FROM users WHERE username ILIKE $1) as exists;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: CreatePlayerResources :exec
INSERT INTO user_resources (user_id)
VALUES ($1);

-- name: GetUserWithResources :one
SELECT u.*, r.*
FROM users u
LEFT JOIN user_resources r ON u.id = r.user_id
WHERE u.id = $1;

-- name: AddUserBits :exec
UPDATE user_resources
SET bits = bits + $2
WHERE user_id = $1;

-- name: GetUserTeam :many
SELECT ut.team_slot, ud.*, bd.*
FROM user_team ut
JOIN user_digimon ud ON ut.digimon_id = ud.id
JOIN base_digimon bd ON ud.base_id = bd.id
WHERE ut.user_id = $1
ORDER BY ut.team_slot;

-- name: AddDigimonToTeam :exec
INSERT INTO user_team (user_id, digimon_id, team_slot)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, team_slot)
DO UPDATE SET digimon_id = EXCLUDED.digimon_id;

-- name: RemoveDigimonFromTeam :exec
DELETE FROM user_team
WHERE user_id = $1 AND digimon_id = $2;

-- name: ClearTeamSlot :exec
DELETE FROM user_team
WHERE user_id = $1 AND team_slot = $2;

-- name: GetUserDigibank :many
SELECT ud.*, bd.*
FROM user_digibank ub
JOIN user_digimon ud ON ub.digimon_id = ud.id
JOIN base_digimon bd ON ud.base_id = bd.id
WHERE ub.user_id = $1
ORDER BY ub.bank_slot;

-- name: CompleteStage :exec
INSERT INTO user_completed_stages (user_id, stage_id)
VALUES ($1, $2)
ON CONFLICT (user_id, stage_id) DO NOTHING;

-- name: GetCompletedStages :many
SELECT * FROM user_completed_stages WHERE user_id = $1;

-- name: GetUserUnlockedChapters :many
SELECT 
    c.*, 
    COALESCE(ucp.is_beaten, false) as is_beaten
FROM chapters c
LEFT JOIN user_chapter_progress ucp ON c.id = ucp.chapter_id AND ucp.user_id = $1
WHERE c.id = 1
   OR COALESCE(ucp.is_unlocked, false) = true
ORDER BY c.chapter_number;

-- name: UnlockUserChapter :exec
INSERT INTO user_chapter_progress (user_id, chapter_id, is_unlocked)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, chapter_id) 
DO UPDATE SET is_unlocked = true;

-- name: CompleteUserChapter :exec
UPDATE user_chapter_progress
SET is_beaten = true, completed_at = NOW()
WHERE user_id = $1 AND chapter_id = $2;

-- name: CheckUserChapterCompletion :one
SELECT COUNT(*) as total_stages,
    COUNT(ucs.stage_id) as completed_stages
FROM stages s
JOIN episodes e ON s.episode_id = e.id
LEFT JOIN user_completed_stages ucs ON s.id = ucs.stage_id AND ucs.user_id = $1
WHERE e.chapter_id = $2
ORDER BY e.chapter_id;

-- name: GetAllUserFlags :many
SELECT * FROM user_flags WHERE user_id = $1;
-- name: GetUserFlagByName :one
SELECT * FROM user_flags WHERE user_id = $1 AND flag_name = $2;
-- name: CreateUserFlag :one
INSERT INTO user_flags (user_id, flag_name, is_active) VALUES ($1, $2, $3) RETURNING *;
-- name: UpdateUserFlag :one
UPDATE user_flags SET is_active = $3 WHERE user_id = $1 AND flag_name = $2 RETURNING *;


-- name: GetAllUserAchievements :many
-- SELECT * FROM user_achievements WHERE user_id = $1;
SELECT uas.*,
    ag.name as achievement_group_name,
    a.name as achievement_name
FROM user_achievements uas
JOIN achievements a ON uas.achievement_id = a.id
JOIN achievement_groups ag on a.catergory_id = ag.id
WHERE uas.user_id = $1;
-- name: GetUserAchievementById :one
SELECT * FROM user_achievements WHERE achievement_id = $1;
-- name: GetUserAchivementsByCategory :many
SELECT ua.*,
    a.name as achievement_name,
    a.category_id,
    ag.name as category_name
FROM user_achievements ua
JOIN achievements a ON ua.achievement_id = a.id
JOIN achievement_groups ag on a.category_id = ag.id
WHERE ua.user_id = $1
    AND ($2 IS NULL OR a.catergory_id = $2)
ORDER BY ag.position, a.name;
-- name: GiveUserAchievement :one
INSERT INTO user_achievements (user_id, achievement_id, current_progress, is_complete)
SELECT $1, a.id, 0, false
FROM achievements a
WHERE a.id = $2
RETURNING *;

-- name: GetAchievementById :one
SELECT * FROM achievements WHERE id = $1;