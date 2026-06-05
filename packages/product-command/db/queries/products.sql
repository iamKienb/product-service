-- name: CreateProduct :exec
INSERT INTO products (
    id,
    shop_id,
    name, 
    slug, 
    description, 
    brand, 
    thumb_url, 
    video_url, 
    price_min, 
    price_max, 
    status, 
    has_variant, 
    created_by,
    created_at
)
VALUES (
    @id::uuid,
    @shop_id::uuid,
    @name::text,
    @slug::text, 
    @description::text, 
    @brand::text, 
    @thumb_url::text, 
    @video_url::text, 
    @price_min::bigint, 
    @price_max::bigint, 
    @status::text, 
    @has_variant::boolean, 
    @created_by::uuid, 
    @created_at::timestamptz
);

-- name: ListProductsByShop :many
SELECT * FROM products
WHERE shop_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CheckProductSlugExists :one
SELECT EXISTS (
    SELECT 1
    FROM products
    WHERE shop_id = @shop_id::uuid
      AND slug = @slug::text
);

-- name: GetProductByID :one
SELECT *
FROM products
WHERE id = @id::uuid
LIMIT 1;

-- name: GetProductByShopAndSlug :one
SELECT *
FROM products
WHERE shop_id = @shop_id::uuid
  AND slug = @slug::text
LIMIT 1;

-- name: DeleteProduct :execrows
DELETE FROM products
WHERE id = @id::uuid;

-- -- name: GetProductDetailWithVariants :one
-- SELECT 
--     p.id,
--     p.name,
--     p.slug,
--     p.description,
--     p.brand,
--     p.thumb_url,
--     p.video_url,
--     p.price_min,
--     p.price_max,
--     p.status,
--     p.has_variant,
--     (
--         SELECT json_agg(json_build_object(
--             'attribute_id', pa.id,
--             'attribute_name', pa.name,
--             'values', (
--                 SELECT json_agg(json_build_object(
--                     'value_id', av.id,
--                     'value_name', av.name
--                 ))
--                 FROM attribute_values av
--                 WHERE av.product_attribute_id = pa.id
--             )
--         ))
--         FROM product_attributes pa
--         WHERE pa.product_id = p.id
--     ) AS attributes_json,
--     (
--         SELECT json_agg(json_build_object(
--             'variant_id', pv.id,
--             'sku_code', pv.sku_code,
--             'price', pv.price,
--             'currency', pv.currency,
--             'image_url', pv.image_url,
--             'status', pv.status,
--             'is_default', pv.is_default,
--             'attribute_value_ids', (
--                 SELECT json_agg(pav.attribute_value_id)
--                 FROM product_attribute_values pav
--                 WHERE pav.attribute_value_id = pv.id
--             )
--         ))
--         FROM product_variants pv
--         WHERE pv.product_id = p.id
--     ) AS variants_json
-- FROM products p
-- WHERE p.id = $1 AND p.shop_id = $2;
