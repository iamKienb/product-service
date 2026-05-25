package product

import (
	"product-command-module/internal/domain/shared"
)

type AttributeParam struct {
	Name   string
	Values []string
}

type VariantParam struct {
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

	HasVariant bool
}
