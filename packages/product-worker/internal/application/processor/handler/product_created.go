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

	doc := map[string]any{
		"product_id":  payload.ProductID,
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

		"attributes": payload.Attributes,
		"variants":   payload.Variants,
	}

	return h.repo.SyncData(ctx, h.alias, payload.ProductID, doc)
}
