package create_product

import (
	"context"
	"product-command-module/internal/domain/shared"
)

type ProductAttribute struct {
	Name   string
	Values []string
}

type ProductVariant struct {
	SkuCode             string
	Price               int64
	Currency            string
	ImageUrl            string
	AttributeValueNames []string
	Quantity            int32
}

type Command struct {
	ShopID      shared.ShopID
	UserID      shared.UserID
	Name        string
	Slug        string
	Description string
	Brand       string
	ThumbUrl    string
	VideoUrl    string
	Attributes  []ProductAttribute
	Variants    []ProductVariant
	HasVariant  bool
	Action      string
}

type SkuItem struct {
	SkuID    string
	Quantity int32
}

type Result struct {
	ProductID string
	ShopID    string
	SkuItems  []SkuItem
}

type Executor interface {
	Execute(ctx context.Context, cmd Command) (*Result, error)
}
