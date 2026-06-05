-- INSERT INTO category_products (
--     shop_id,
--     category_id,
--     product_id
-- )
-- FROM (
--     SELECT * FROM ROWS FROM(
--         unnest(@shop_id::uuid[]),
--         unnest(@category_id::uuid[]), 
--         unnest(@product_id::uuid[])
--     ) AS tmp(shop_id, category_id, product_id)
-- ) AS u

-- ON CONFLICT (shop_id, category_id, product_id) DO NOTHING;

-- name: RemoveProductFromCategory :exec
DELETE FROM category_products
WHERE shop_id = @shop_id::uuid 
    AND category_id = @category_id::uuid 
    AND product_id = @product_id::uuid;

-- name: RemoveProductFromAllCategories :exec
DELETE FROM category_products
WHERE shop_id = @shop_id::uuid 
    AND product_id = @product_id::uuid;

-- name: ListCategoriesByProductID :many
SELECT c.*
FROM categories c
JOIN category_products cp ON c.id = cp.category_id
WHERE cp.shop_id = $1 AND cp.product_id = $2;

-- name: CountBySlug :one
SELECT
    1
FROM categories 
WHERE slug = @slug::text;
