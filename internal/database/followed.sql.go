// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: followed.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFollow = `-- name: CreateFollow :one
INSERT INTO follow_users(id,created_at,account_name,hashtag,puuID,region)
VALUES($1,$2,$3,$4,$5,$6)
RETURNING id, created_at, account_name, hashtag, puuid, region
`

type CreateFollowParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	AccountName string
	Hashtag     string
	Puuid       string
	Region      string
}

func (q *Queries) CreateFollow(ctx context.Context, arg CreateFollowParams) (FollowUser, error) {
	row := q.db.QueryRowContext(ctx, createFollow,
		arg.ID,
		arg.CreatedAt,
		arg.AccountName,
		arg.Hashtag,
		arg.Puuid,
		arg.Region,
	)
	var i FollowUser
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.AccountName,
		&i.Hashtag,
		&i.Puuid,
		&i.Region,
	)
	return i, err
}

const getFollowed = `-- name: GetFollowed :many
Select id, created_at, account_name, hashtag, puuid, region FROM follow_users
`

func (q *Queries) GetFollowed(ctx context.Context) ([]FollowUser, error) {
	rows, err := q.db.QueryContext(ctx, getFollowed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FollowUser
	for rows.Next() {
		var i FollowUser
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.AccountName,
			&i.Hashtag,
			&i.Puuid,
			&i.Region,
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
