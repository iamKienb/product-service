-- name: CreateCategory :exec
INSERT INTO categories (
    id, 
    shop_id, 
    status, 
    parent_id, 
    name, 
    slug, 
    created_by,
    created_at
) 
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: GetCategoryByID :one
SELECT * FROM categories 
WHERE id = $1 AND shop_id = $2;

-- name: GetCategoryBySlug :one
SELECT * FROM categories 
WHERE shop_id = $1 AND slug = $2;

-- name: UpdateCategory :one
UPDATE categories
SET 
    name = $1,
    slug = $2,
    status = $3,
    parent_id = $4,
    updated_by = $5,
    updated_at = NOW()
WHERE id = $6 AND shop_id = $7
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories 
WHERE id = $1 AND shop_id = $2;

-- name: GetCategoryChildren :many
SELECT * FROM categories
WHERE shop_id = $1 AND parent_id = $2
ORDER BY name ASC;
