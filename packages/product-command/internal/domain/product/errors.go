package product

import "errors"

var (
	ErrInvalidProductID        = errors.New("invalid_product_id")
	ErrEmptyName               = errors.New("name must not be empty")
	ErrEmptySKUCode            = errors.New("sku code must not be empty")
	ErrDuplicateSKUCode        = errors.New("duplicate sku code")
	ErrNoSKU                   = errors.New("must have at least one sku")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrNegativePrice           = errors.New("price must not be negative")
	ErrEmptyCurrency           = errors.New("currency must not be empty")
	ErrProductSlugTaken        = errors.New("product_slug_taken")
)
