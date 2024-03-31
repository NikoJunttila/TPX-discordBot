-- name: CreateFollow :one
INSERT INTO follow_users(id,created_at,account_name,hashtag,puuID,region)
VALUES(?,?,?,?,?,?)
RETURNING *;
--
-- name: GetFollowed :many
Select * FROM follow_users;

