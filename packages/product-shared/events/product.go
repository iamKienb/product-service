package events

import "time"

const (
	TopicProductCreated = "product-service.product.created"
	TopicProductDeleted = "product-service.product.deleted"
)

type ProductCreated struct {
	ProductID   string           `json:"product_id"`
	ShopID      string           `json:"shop_id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Description string           `json:"description"`
	Brand       string           `json:"brand"`
	ThumbURL    string           `json:"thumb_url"`
	VideoURL    string           `json:"video_url"`
	PriceMin    int64            `json:"price_min"`
	PriceMax    int64            `json:"price_max"`
	Status      string           `json:"status"`
	HasVariant  bool             `json:"has_variant"`
	CreatedBy   string           `json:"created_by"`
	CreatedAt   time.Time        `json:"created_at"`
	Attributes  []AttributeEvent `json:"attributes"` // Gom cụm attribute
	Variants    []VariantEvent   `json:"variants"`   // Gom cụm variant
}

type AttributeEvent struct {
	AttributeID   string                `json:"attribute_id"`
	AttributeName string                `json:"attribute_name"`
	Values        []AttributeValueEvent `json:"values"`
}

type AttributeValueEvent struct {
	AttributeValueID string `json:"attribute_value_id"`
	ValueName        string `json:"value_name"`
}

type VariantEvent struct {
	SkuID             string   `json:"sku_id"`
	SkuCode           string   `json:"sku_code"`
	Price             int64    `json:"price"`
	Currency          string   `json:"currency"`
	ImageURL          string   `json:"image_url"`
	Status            string   `json:"status"`
	IsDefault         bool     `json:"is_default"`
	AttributeValueIDs []string `json:"attribute_value_ids"`
}

type ProductDeleted struct {
	ProductID string `json:"product_id"`
	ShopID    string `json:"shop_id"`
}
