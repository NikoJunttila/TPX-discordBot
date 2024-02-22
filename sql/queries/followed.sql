-- name: CreateFollow :one
INSERT INTO follow_users(id,created_at,account_name,hashtag,puuID,region)
VALUES($1,$2,$3,$4,$5,$6)
RETURNING *;
--
-- name: GetFollowed :many
Select * FROM follow_users;

