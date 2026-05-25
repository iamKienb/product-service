CREATE Table products (
    id UUID PRIMARY KEY,
    shop_id UUID NOT NULL,

    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT NOT NULL,
    brand TEXT NOT NULL,

    thumb_url TEXT NOT NULL,
    video_url TEXT NOT NULL,

    price_min BIGINT NOT NULL, 
    price_max BIGINT NOT NULL,
    status TEXT NOT NULL,
    has_variant BOOLEAN DEFAULT FALSE,

    created_by UUID NOT NULL,
    updated_by UUID DEFAULT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NULL,

    UNIQUE(shop_id, slug),
    CONSTRAINT product_status_check
        CHECK (status IN ('DRAFT', 'ACTIVE', 'ARCHIVED'))
);

CREATE INDEX idx_products_shop_id ON products(shop_id);
CREATE INDEX idx_products_slug ON products(slug);
