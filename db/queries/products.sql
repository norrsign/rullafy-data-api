-- name: CreateProduct :one
INSERT INTO products (id, name, type, description)
VALUES ($1, $2, $3, $4)
RETURNING id, name, type, description;

-- name: GetProduct :one
SELECT id, name, type, description
FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT id, name, type, description
FROM products
ORDER BY name;

-- name: UpdateProduct :one
UPDATE products
SET name = $2,
    type = $3,
    description = $4
WHERE id = $1
RETURNING id, name, type, description;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;
