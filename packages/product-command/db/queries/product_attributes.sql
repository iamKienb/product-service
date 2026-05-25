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


-- name: BatchCreateAttributeValues :exec
INSERT INTO attribute_values (
    id, 
    product_attribute_id, 
    value
) 
SELECT
    unnest(@ids::uuid[]),
    unnest(@product_attribute_ids::uuid[]),
    unnest(@values::text[]);
