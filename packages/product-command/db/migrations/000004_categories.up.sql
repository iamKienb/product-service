CREATE TABLE categories(
    id UUID PRIMARY KEY,
    shop_id UUID NOT NULL,
    status TEXT NOT NULL,
    parent_id UUID DEFAULT NULL REFERENCES categories(id),
    name TEXT NOT NULL,
    slug TEXT NOT NULL,

    created_by UUID NOT NULL,
    updated_by UUID DEFAULT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NULL,
    UNIQUE(shop_id, slug),

    CONSTRAINT category_status_check
        CHECK (status IN ('DRAFT', 'ACTIVE', 'ARCHIVED'))
);

CREATE TABLE category_products(
    shop_id UUID NOT NULL,
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    PRIMARY KEY (shop_id, category_id, product_id)
);

CREATE INDEX idx_categories_shop_id ON categories(shop_id);
CREATE INDEX idx_categories_slug ON categories(slug);
