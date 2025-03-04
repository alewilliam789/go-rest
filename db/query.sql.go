// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :execresult
INSERT INTO user (
  user_name,
  password,
  first_name,
  last_name,
  dob,
  city,
  state
) values (
  ?,?,?,?,?,?,?
)
`

type CreateUserParams struct {
	UserName  string
	Password  string
	FirstName sql.NullString
	LastName  sql.NullString
	Dob       sql.NullString
	City      sql.NullString
	State     sql.NullString
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createUser,
		arg.UserName,
		arg.Password,
		arg.FirstName,
		arg.LastName,
		arg.Dob,
		arg.City,
		arg.State,
	)
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM user
WHERE user_id=?
`

func (q *Queries) DeleteUser(ctx context.Context, userID int32) error {
	_, err := q.db.ExecContext(ctx, deleteUser, userID)
	return err
}

const getUser = `-- name: GetUser :one
SELECT user_id, user_name, password, first_name, last_name, dob, city, state FROM user
WHERE user_id = ? LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, userID int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.UserName,
		&i.Password,
		&i.FirstName,
		&i.LastName,
		&i.Dob,
		&i.City,
		&i.State,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE user
SET password=?, first_name=?,last_name=?,dob=?,city=?,state=?
WHERE user_id = ?
`

type UpdateUserParams struct {
	Password  string
	FirstName sql.NullString
	LastName  sql.NullString
	Dob       sql.NullString
	City      sql.NullString
	State     sql.NullString
	UserID    int32
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.Password,
		arg.FirstName,
		arg.LastName,
		arg.Dob,
		arg.City,
		arg.State,
		arg.UserID,
	)
	return err
}
