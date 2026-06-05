package product

import (
	"product-command-module/internal/domain/shared"
	"time"
)

type AttributeParam struct {
	Name   string
	Values []string
}

type ProductVariantParams struct {
	SkuCode             string
	Price               int64
	Currency            string
	ImageUrl            string
	AttributeValueNames []string
}

type NewProductParams struct {
	ShopID      shared.ShopID
	UserID      shared.UserID
	Name        string
	Slug        string
	Description string
	Brand       string
	ThumbUrl    string
	VideoUrl    string
	HasVariant  bool
}

type NewVariantParams struct {
	SkuID             shared.SkuID
	ProductID         shared.ProductID
	ShopID            shared.ShopID
	SkuCode           string
	Price             int64
	Currency          string
	ImageUrl          string
	IsDefault         bool
	AttributeValueIDs []shared.AttributeValueID
	CreatedBy         shared.UserID
	CreatedAt         time.Time
}
