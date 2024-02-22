-- name: CreateUser :one
INSERT INTO usersCount(id,created_at,updated_at,count)
VALUES($1,$2,$3,$4)
RETURNING *;
--
-- name: GetUser :one
Select * FROM usersCount WHERE id = $1;
-- name: UpdateUser :exec
UPDATE usersCount
  set count = count + $2
WHERE id = $1;
--
-- name: HighscoreUsers :many
SELECT *
FROM userscount
ORDER BY count DESC
LIMIT 5;
--