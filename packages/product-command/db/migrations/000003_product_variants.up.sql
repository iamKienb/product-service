CREATE Table product_variants (
    sku_id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    shop_id UUID NOT NULL,
    sku_code TEXT NOT NULL,
    price BIGINT NOT NULL,
    currency TEXT NOT NULL DEFAULT 'VND',
    image_url TEXT NOT NULL,
    status TEXT NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,

    created_by UUID NOT NULL,
    updated_by UUID DEFAULT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NULL,

    UNIQUE(shop_id, sku_code),
    CONSTRAINT product_variant_status_check
        CHECK (status IN('IN_STOCK', 'OUT_OF_STOCK', 'COMING_SOON', 'DISCONTINUED', 'HIDDEN'))
);

CREATE Table product_attribute_values (
    sku_id UUID NOT NULL REFERENCES product_variants(sku_id) ON DELETE CASCADE,
    attribute_value_id UUID NOT NULL REFERENCES attribute_values(id) ON DELETE CASCADE,
    primary key (sku_id, attribute_value_id)
);

CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);
