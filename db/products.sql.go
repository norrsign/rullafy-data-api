// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: products.sql

package db

import (
	"context"
)

const createProduct = `-- name: CreateProduct :one
INSERT INTO products (id, name, type, description)
VALUES ($1, $2, $3, $4)
RETURNING id, name, type, description
`

type CreateProductParams struct {
	ID          string
	Name        string
	Type        string
	Description string
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, createProduct,
		arg.ID,
		arg.Name,
		arg.Type,
		arg.Description,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.Description,
	)
	return i, err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteProduct, id)
	return err
}

const getProduct = `-- name: GetProduct :one
SELECT id, name, type, description
FROM products
WHERE id = $1
`

func (q *Queries) GetProduct(ctx context.Context, id string) (Product, error) {
	row := q.db.QueryRow(ctx, getProduct, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.Description,
	)
	return i, err
}

const listProducts = `-- name: ListProducts :many
SELECT id, name, type, description
FROM products
ORDER BY name
`

func (q *Queries) ListProducts(ctx context.Context) ([]Product, error) {
	rows, err := q.db.Query(ctx, listProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Type,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateProduct = `-- name: UpdateProduct :one
UPDATE products
SET name = $2,
    type = $3,
    description = $4
WHERE id = $1
RETURNING id, name, type, description
`

type UpdateProductParams struct {
	ID          string
	Name        string
	Type        string
	Description string
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, updateProduct,
		arg.ID,
		arg.Name,
		arg.Type,
		arg.Description,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.Description,
	)
	return i, err
}
