-- name: GetUser :one
SELECT * FROM user
WHERE user_name = ? LIMIT 1;

-- name: CreateUser :execresult
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
);

-- name: UpdateUser :exec
UPDATE user
SET password=?, first_name=?,last_name=?,dob=?,city=?,state=?
WHERE user_name = ?;

-- name: DeleteUser :exec
DELETE FROM user
WHERE user_name=?;

