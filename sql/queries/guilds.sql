-- name: CreateGuild :one
INSERT INTO guildCount(id,created_at,updated_at,count)
VALUES(?,?,?,?)
RETURNING *;
--
-- name: GetGuild :one
Select * FROM guildCount WHERE id = ?;
--
-- name: UpdateGuild :exec
UPDATE guildCount
  set count = count + ?
WHERE id = ?;
--
-- name: HighscoreGuild :many
SELECT *
FROM guildCount
ORDER BY count DESC
LIMIT 5;
--
