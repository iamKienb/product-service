-- name: BatchCreateAttributes :exec
INSERT INTO product_attributes (
    id,
    product_id,
    name
)
SELECT
    unnest(@ids::uuid[]),
    @product_id::uuid,
    unnest(@names::text[]);

-- name: ListAttributesByProductID :many
SELECT *
FROM product_attributes
WHERE product_id = @product_id::uuid
ORDER BY name ASC;

-- name: BatchCreateAttributeValues :exec
INSERT INTO attribute_values (
    id, 
    product_attribute_id, 
    name
) 
SELECT
    unnest(@ids::uuid[]),
    unnest(@product_attribute_ids::uuid[]),
    unnest(@values::text[]);

-- name: ListAttributeValuesByProductID :many
SELECT av.*
FROM attribute_values av
JOIN product_attributes pa ON pa.id = av.product_attribute_id
WHERE pa.product_id = @product_id::uuid
ORDER BY av.name ASC;
