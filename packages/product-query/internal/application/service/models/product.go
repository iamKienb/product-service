package models

import "time"

type Product struct {
	ID          string           `json:"id"`
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
	Variants    []ProductVariant `json:"variants"`
	Attributes  []ProductAttr    `json:"attributes"`

	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductVariant struct {
	SkuID             string   `json:"id"`
	SkuCode           string   `json:"code"`
	Price             int64    `json:"price"`
	Currency          string   `json:"currency"`
	ImageURL          string   `json:"image_url"`
	Status            string   `json:"status"`
	IsDefault         bool     `json:"is_default"`
	AttributeValueIDs []string `json:"attribute_value_ids"`

	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Value struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type ProductAttr struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Values []Value `json:"values"`
}
