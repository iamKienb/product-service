-- name: BatchCreateVariants :exec
INSERT INTO product_variants (
    sku_id,
    product_id, 
    shop_id,
    sku_code,
    price,
    currency,
    image_url,
    status,
    is_default,
    created_by,
    created_at
)
SELECT
    unnest(@sku_ids::uuid[]),
    @product_id::uuid,
    @shop_id::uuid,
    unnest(@sku_codes::text[]),
    unnest(@prices::bigint[]),
    unnest(@currencies::text[]),
    unnest(@image_urls::text[]),
    unnest(@status::text[]),
    unnest(@is_defaults::boolean[]),
    unnest(@created_bys::uuid[]),
    unnest(@created_ats::timestamptz[]);

-- name: BatchLinkVariantAttributes :exec
INSERT INTO product_attribute_values (
    sku_id,
    attribute_value_id
)
SELECT
    unnest(@sku_ids::uuid[]),
    unnest(@attribute_value_ids::uuid[])
ON CONFLICT DO NOTHING;


-- name: FindPriceSkusByIDs :many
SELECT 
    sku_id, 
    price
FROM product_variants
WHERE sku_id = ANY(@sku_ids::uuid[])
    AND shop_id = @shop_id::uuid;

-- name: ListVariantsByProductID :many
SELECT *
FROM product_variants
WHERE product_id = @product_id::uuid;

-- name: ListVariantAttributeValuesByProductID :many
SELECT pav.*
FROM product_attribute_values pav
JOIN product_variants pv ON pv.sku_id = pav.sku_id
WHERE pv.product_id = @product_id::uuid;
