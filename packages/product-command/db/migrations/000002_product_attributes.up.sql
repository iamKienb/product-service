CREATE Table product_attributes(
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    UNIQUE(product_id, name)
);

CREATE TABLE attribute_values(
    id UUID PRIMARY KEY,
    product_attribute_id UUID NOT NULL REFERENCES product_attributes(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    UNIQUE(product_attribute_id, name)
);