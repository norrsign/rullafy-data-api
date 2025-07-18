-- name: CreateCompany :one
INSERT INTO companies (id, name, address, phone)
VALUES ($1, $2, $3, $4)
RETURNING id, name, address, phone;

-- name: GetCompany :one
SELECT id, name, address, phone
FROM companies
WHERE id = $1;

-- name: ListCompanies :many
SELECT id, name, address, phone
FROM companies
ORDER BY name;

-- name: UpdateCompany :one
UPDATE companies
SET name = $2,
    address = $3,
    phone = $4
WHERE id = $1
RETURNING id, name, address, phone;

-- name: DeleteCompany :exec
DELETE FROM companies
WHERE id = $1;
