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

		"thumb_url": e.ThumbUrl,
		"video_url": e.VideoUrl,

		"price_min": e.PriceMin,
		"price_max": e.PriceMax,

		"status":     string(e.Status),
		"hasVariant": e.HasVariant,

		"created_by": e.CreatedBy.String(),
		"created_at": e.CreatedAt,
	}
}

type VariantCreatedEvent struct {
	SkuID             shared.SkuID
	ProductID         shared.ProductID
	ShopID            shared.ShopID
	SkuCode           string
	Price             int64
	Currency          string
	ImageUrl          string
	AttributeValueIDs []shared.AttributeValueID

	CreatedBy shared.UserID
	CreatedAt time.Time
}

func (e VariantCreatedEvent) EventName() string {
	return "product-service.variant.created"
}

func (e VariantCreatedEvent) IntegrationPayload() map[string]interface{} {
	attrValueStrings := make([]string, 0, len(e.AttributeValueIDs))

	for _, valueID := range e.AttributeValueIDs {
		attrValueStrings = append(attrValueStrings, valueID.String())
	}

	return map[string]interface{}{
		"sku_id":              e.SkuID.String(),
		"product_id":          e.ProductID.String(),
		"shop_id":             e.ShopID.String(),
		"sku_code":            e.SkuCode,
		"currency":            e.Currency,
		"ImageUrl":            e.ImageUrl,
		"attribute_value_ids": attrValueStrings,

		"created_by": e.CreatedBy.String(),
		"created_at": e.CreatedAt,
	}
}

type AttributeCreatedEvent struct {
	AttributeID   shared.AttributeID
	ProductID     shared.ProductID
	AttributeName string
}

func (e AttributeCreatedEvent) EventName() string {
	return "product-service.attribute.created"
}

func (e AttributeCreatedEvent) IntegrationPayload() map[string]interface{} {
	return map[string]interface{}{
		"attribute_id":   e.AttributeID.String(),
		"product_id":     e.ProductID.String(),
		"attribute_name": e.AttributeName,
	}
}

type AttributeValueCreatedEvent struct {
	AttributeValueID shared.AttributeValueID
	AttributeID      shared.AttributeID
	ValueName        string
}

func (e AttributeValueCreatedEvent) EventName() string {
	return "product-service.attribute.value.created"
}

func (e AttributeValueCreatedEvent) IntegrationPayload() map[string]interface{} {
	return map[string]interface{}{
		"attribute_value_id": e.AttributeValueID.String(),
		"attribute_id":       e.AttributeID.String(),
		"value_name":         e.ValueName,
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
