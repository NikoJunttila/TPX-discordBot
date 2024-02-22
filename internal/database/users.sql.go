// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: users.sql

package database

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO usersCount(id,created_at,updated_at,count)
VALUES($1,$2,$3,$4)
RETURNING id, created_at, updated_at, count
`

type CreateUserParams struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Count     int32
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (Userscount, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Count,
	)
	var i Userscount
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Count,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
Select id, created_at, updated_at, count FROM usersCount WHERE id = $1
`

func (q *Queries) GetUser(ctx context.Context, id string) (Userscount, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i Userscount
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Count,
	)
	return i, err
}

const highscoreUsers = `-- name: HighscoreUsers :many
SELECT id, created_at, updated_at, count
FROM userscount
ORDER BY count DESC
LIMIT 5
`

func (q *Queries) HighscoreUsers(ctx context.Context) ([]Userscount, error) {
	rows, err := q.db.QueryContext(ctx, highscoreUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Userscount
	for rows.Next() {
		var i Userscount
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Count,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :exec
UPDATE usersCount
  set count = count + $2
WHERE id = $1
`

type UpdateUserParams struct {
	ID    string
	Count int32
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser, arg.ID, arg.Count)
	return err
}
