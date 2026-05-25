package product

import (
	"product-command-module/internal/domain/shared"
	"time"
)

type Variant struct {
	SkuID             shared.SkuID
	ProductID         shared.ProductID
	ShopID            shared.ShopID
	SkuCode           string
	Price             int64
	Currency          string
	ImageUrl          string
	IsDefault         bool
	AttributeValueIDs []shared.AttributeValueID

	CreatedBy shared.UserID
	UpdatedBy *shared.UserID

	CreatedAt time.Time
	UpdatedAt *time.Time
}
