-- name: CreateGuild :one
INSERT INTO guildCount(id,created_at,updated_at,count)
VALUES($1,$2,$3,$4)
RETURNING *;
--
-- name: GetGuild :one
Select * FROM guildCount WHERE id = $1;
--
--
UPDATE guildCount
  set count = count + $2,
WHERE id = $1;
--