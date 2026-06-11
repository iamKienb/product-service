package handler

import (
	"context"
	"encoding/json"
	"product-shared-module/events"
	"product-worker-module/internal/application/port"
)

type ProductCreatedHandler struct {
	repo  port.ESRepository
	alias string
}

func NewProductCreatedHandler(repo port.ESRepository, alias string) *ProductCreatedHandler {
	return &ProductCreatedHandler{repo: repo, alias: alias}
}

func (h *ProductCreatedHandler) Handle(ctx context.Context, raw json.RawMessage) error {
	var payload events.ProductCreated
	if err := json.Unmarshal(raw, &payload); err != nil {
		return err
	}

	esAttributes := make([]map[string]any, 0, len(payload.Attributes))
	for _, attr := range payload.Attributes {
		esValues := make([]map[string]any, 0, len(attr.Values))
		for _, val := range attr.Values {
			esValues = append(esValues, map[string]any{
				"id":   val.AttributeValueID,
				"name": val.ValueName,
			})
		}

		esAttributes = append(esAttributes, map[string]any{
			"id":     attr.AttributeID,
			"name":   attr.AttributeName,
			"values": esValues,
		})
	}

	esVariants := make([]map[string]any, 0, len(payload.Variants))
	for _, variant := range payload.Variants {
		esVariants = append(esVariants, map[string]any{
			"id":                  variant.SkuID,
			"code":                variant.SkuCode,
			"price":               variant.Price,
			"currency":            variant.Currency,
			"image_url":           variant.ImageURL,
			"status":              variant.Status,
			"is_default":          variant.IsDefault,
			"attribute_value_ids": variant.AttributeValueIDs,

			"created_by": payload.CreatedBy,
			"updated_by": payload.CreatedBy,
			"created_at": payload.CreatedAt,
			"updated_at": payload.CreatedAt,
		})
	}

	doc := map[string]any{
		"id":          payload.ProductID,
		"shop_id":     payload.ShopID,
		"name":        payload.Name,
		"slug":        payload.Slug,
		"description": payload.Description,
		"brand":       payload.Brand,
		"thumb_url":   payload.ThumbURL,
		"video_url":   payload.VideoURL,
		"price_min":   payload.PriceMin,
		"price_max":   payload.PriceMax,
		"status":      payload.Status,
		"has_variant": payload.HasVariant,
		"created_by":  payload.CreatedBy,
		"created_at":  payload.CreatedAt,

		"attributes": esAttributes,
		"variants":   esVariants,
	}

	return h.repo.SyncData(ctx, h.alias, payload.ProductID, doc)
}
