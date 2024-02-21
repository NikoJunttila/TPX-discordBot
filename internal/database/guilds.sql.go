// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: guilds.sql

package database

import (
	"context"
	"time"
)

const createGuild = `-- name: CreateGuild :one
INSERT INTO guildCount(id,created_at,updated_at,count)
VALUES($1,$2,$3,$4)
RETURNING id, created_at, updated_at, count
`

type CreateGuildParams struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Count     int32
}

func (q *Queries) CreateGuild(ctx context.Context, arg CreateGuildParams) (Guildcount, error) {
	row := q.db.QueryRowContext(ctx, createGuild,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Count,
	)
	var i Guildcount
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Count,
	)
	return i, err
}

const getGuild = `-- name: GetGuild :one
Select id, created_at, updated_at, count FROM guildCount WHERE id = $1
`

func (q *Queries) GetGuild(ctx context.Context, id string) (Guildcount, error) {
	row := q.db.QueryRowContext(ctx, getGuild, id)
	var i Guildcount
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Count,
	)
	return i, err
}

const updateGuild = `-- name: UpdateGuild :exec
UPDATE guildCount
  set count = count + $2
WHERE id = $1
`

type UpdateGuildParams struct {
	ID    string
	Count int32
}

func (q *Queries) UpdateGuild(ctx context.Context, arg UpdateGuildParams) error {
	_, err := q.db.ExecContext(ctx, updateGuild, arg.ID, arg.Count)
	return err
}
