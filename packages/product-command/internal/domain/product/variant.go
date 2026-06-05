package product

import (
	"product-command-module/internal/domain/shared"
	"time"
)

type ProductVariant struct {
	SkuID             shared.SkuID
	ProductID         shared.ProductID
	ShopID            shared.ShopID
	SkuCode           string
	Price             int64
	Currency          string
	ImageUrl          string
	Status            VariantStatus
	IsDefault         bool
	AttributeValueIDs []shared.AttributeValueID

	CreatedBy shared.UserID
	UpdatedBy *shared.UserID

	CreatedAt time.Time
	UpdatedAt *time.Time
}

type VariantStatus string

const (
	VariantStatusActive   VariantStatus = "ACTIVE"
	VariantStatusInactive VariantStatus = "INACTIVE"
	VariantStatusArchived VariantStatus = "ARCHIVED"
)

func NewProductVariant(params NewVariantParams) *ProductVariant {
	return &ProductVariant{
		SkuID:             params.SkuID,
		ProductID:         params.ProductID,
		ShopID:            params.ShopID,
		SkuCode:           params.SkuCode,
		Price:             params.Price,
		Currency:          params.Currency,
		ImageUrl:          params.ImageUrl,
		Status:            VariantStatusActive,
		IsDefault:         params.IsDefault,
		AttributeValueIDs: params.AttributeValueIDs,

		CreatedBy: params.CreatedBy,
		CreatedAt: params.CreatedAt,

		UpdatedAt: nil,
		UpdatedBy: nil,
	}
}
