INSERT INTO category_products (
    shop_id, 
    category_id, 
    product_id
)
VALUES ($1, $2, $3)
ON CONFLICT (shop_id, category_id, product_id) DO NOTHING;

-- name: RemoveProductFromCategory :exec
DELETE FROM category_products
WHERE shop_id = $1 AND category_id = $2 AND product_id = $3;

-- name: RemoveProductFromAllCategories :exec
DELETE FROM category_products
WHERE shop_id = $1 AND product_id = $2;

-- name: ListCategoriesByProductID :many
SELECT c.* 
FROM categories c
JOIN category_products cp ON c.id = cp.category_id
WHERE cp.shop_id = $1 AND cp.product_id = $2;

-- name: CountBySlug :one
SELECT 
    1
FROM categories 
WHERE slug = $1;
