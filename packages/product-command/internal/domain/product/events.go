package product

import (
	"product-command-module/internal/domain/shared"
	"time"
)

type ProductCreatedEvent struct {
	ProductID   shared.ProductID
	ShopID      shared.ShopID
	Name        string
	Slug        string
	Description string
	Brand       string
	ThumbUrl    string
	VideoUrl    string

	PriceMin int64
	PriceMax int64

	Status     ProductStatus
	HasVariant bool

	CreatedBy shared.UserID
	CreatedAt time.Time

	Attributes []AttributePayload
	Variants   []VariantPayload
}

type AttributePayload struct {
	AttributeID   string                  `json:"attribute_id"`
	AttributeName string                  `json:"attribute_name"`
	Values        []AttributeValuePayload `json:"values"`
}

type AttributeValuePayload struct {
	AttributeValueID string `json:"attribute_value_id"`
	ValueName        string `json:"value_name"`
}

type VariantPayload struct {
	SkuID             string   `json:"sku_id"`
	SkuCode           string   `json:"sku_code"`
	Price             int64    `json:"price"`
	Currency          string   `json:"currency"`
	ImageUrl          string   `json:"image_url"`
	Status            string   `json:"status"`
	IsDefault         bool     `json:"is_default"`
	AttributeValueIDs []string `json:"attribute_value_ids"`
}

func (e ProductCreatedEvent) EventName() string {
	return "product-service.product.created"
}

func (e ProductCreatedEvent) IntegrationPayload() map[string]interface{} {
	return map[string]interface{}{
		"product_id":  e.ProductID.String(),
		"shop_id":     e.ShopID.String(),
		"name":        e.Name,
		"slug":        e.Slug,
		"description": e.Description,
		"brand":       e.Brand,
		"thumb_url":   e.ThumbUrl,
		"video_url":   e.VideoUrl,
		"price_min":   e.PriceMin,
		"price_max":   e.PriceMax,
		"status":      string(e.Status),
		"has_variant": e.HasVariant,
		"created_by":  e.CreatedBy.String(),
		"created_at":  e.CreatedAt,

		"attributes": e.Attributes,
		"variants":   e.Variants,
	}
}

type ProductDeletedEvent struct {
	ProductID shared.ProductID
	ShopID    shared.ShopID
}

func (e ProductDeletedEvent) EventName() string {
	return "product-service.product.deleted"
}

func (e ProductDeletedEvent) IntegrationPayload() map[string]interface{} {
	return map[string]interface{}{
		"product_id": e.ProductID.String(),
		"shop_id":    e.ShopID.String(),
	}
}
