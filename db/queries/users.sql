-- name: GetUser :one
-- id: string
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: ListUsers :many
-- limit:  int32
-- offset: int32
SELECT *
FROM users
ORDER BY id
LIMIT  $1
OFFSET $2;

-- name: CreateUser :one
-- id:  string
-- job: string
INSERT INTO users (
  id,
  job
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateUser :one
-- id:  string
-- job: string
UPDATE users
SET job = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
-- id: string
DELETE FROM users
WHERE id = $1;

-- name: CountUsers :one
SELECT COUNT(*)::int4 FROM users;
