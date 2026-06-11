package handler

import (
	"context"
	"encoding/json"
	"product-shared-module/events"
	"product-worker-module/internal/application/port"
)

type ProductDeletedHandler struct {
	repo  port.ESRepository
	alias string
}

func NewProductDeletedHandler(repo port.ESRepository, alias string) *ProductDeletedHandler {
	return &ProductDeletedHandler{repo: repo, alias: alias}
}

func (h *ProductDeletedHandler) Handle(ctx context.Context, raw json.RawMessage) error {
	var payload events.ProductDeleted
	if err := json.Unmarshal(raw, &payload); err != nil {
		return err
	}
	return h.repo.Delete(ctx, h.alias, payload.ProductID)
}
