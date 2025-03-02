-- name: GetUser :one
SELECT * FROM user
WHERE id = ? LIMIT 1;

-- name: CreateUser :execresult
INSERT INTO user (
  username,
  password,
  firstname,
  lastname,
  dob,
  city,
  state
) values (
  ?,?,?,?,?,?,?
);

-- name: UpdateUser :exec
UPDATE user
SET password=?, firstname=?,lastname=?,dob=?,city=?,state=?
WHERE id = ?;

-- name: DeleteAuthor :exec
DELETE FROM user
WHERE id=?;

