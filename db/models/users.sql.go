// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: users.sql

package models

import (
	"context"
)

const countUsers = `-- name: CountUsers :one
SELECT COUNT(*)::int4 FROM users
`

func (q *Queries) CountUsers(ctx context.Context) (int32, error) {
	row := q.db.QueryRow(ctx, countUsers)
	var column_1 int32
	err := row.Scan(&column_1)
	return column_1, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  id,
  job,
  address
) VALUES (
  $1, $2, $3
)
RETURNING id, job, address
`

type CreateUserParams struct {
	ID      string
	Job     string
	Address AddressList
}

// id:      string
// job:     string
// address: models.AddressList
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.ID, arg.Job, arg.Address)
	var i User
	err := row.Scan(&i.ID, &i.Job, &i.Address)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

// id: string
func (q *Queries) DeleteUser(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, job, address
FROM users
WHERE id = $1
LIMIT 1
`

// id: string
func (q *Queries) GetUser(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(&i.ID, &i.Job, &i.Address)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, job, address
FROM users
ORDER BY id
LIMIT  $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32
	Offset int32
}

// limit:  int32
// offset: int32
func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.ID, &i.Job, &i.Address); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchUsersByAddress = `-- name: SearchUsersByAddress :many
SELECT id, job, address
FROM users
WHERE address @> $1  -- JSONB “contains” operator
ORDER BY id
LIMIT  $2
OFFSET $3
`

type SearchUsersByAddressParams struct {
	Address AddressList
	Limit   int32
	Offset  int32
}

// address: AddressList
// limit:  int32
// offset: int32
func (q *Queries) SearchUsersByAddress(ctx context.Context, arg SearchUsersByAddressParams) ([]User, error) {
	rows, err := q.db.Query(ctx, searchUsersByAddress, arg.Address, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.ID, &i.Job, &i.Address); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET job = $2,
    address = $3
WHERE id = $1
RETURNING id, job, address
`

type UpdateUserParams struct {
	ID      string
	Job     string
	Address AddressList
}

// id:      string
// job:     string
// address: models.AddressList
func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser, arg.ID, arg.Job, arg.Address)
	var i User
	err := row.Scan(&i.ID, &i.Job, &i.Address)
	return i, err
}
