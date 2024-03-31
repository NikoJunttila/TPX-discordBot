-- name: CreateUser :one
INSERT INTO usersCount(id,created_at,updated_at,count)
VALUES(?,?,?,?)
RETURNING *;
--
-- name: GetUser :one
Select * FROM usersCount WHERE id = ?;
-- name: UpdateUser :exec
UPDATE usersCount
  set count = count + ?
WHERE id = ?;
--
-- name: HighscoreUsers :many
SELECT *
FROM usersCount
ORDER BY count DESC
LIMIT 5;
--