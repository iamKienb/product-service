package product

import (
	"errors"

	domain_product "product-command-module/internal/domain/product"

	"github.com/iamKienb/go-core/app_error"
)

const (
	errCodeProductSlugTaken      = "product_slug_taken"
	errCodeProductInvalid        = "product_invalid"
	errCodeProductNotFound       = "product_not_found"
	errCodeProductVariantInvalid = "product_variant_invalid"

	errMsgProductSlugTaken      = "product slug is already taken"
	errMsgProductInvalid        = "product data is invalid"
	errMsgProductNotFound       = "product was not found"
	errMsgProductVariantInvalid = "product variant data is invalid"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, domain_product.ErrProductSlugTaken):
		return app_error.New(app_error.KindConflict, errCodeProductSlugTaken, errMsgProductSlugTaken, err)
	case errors.Is(err, domain_product.ErrProductNotFound):
		return app_error.New(app_error.KindNotFound, errCodeProductNotFound, errMsgProductNotFound, err)
	case errors.Is(err, domain_product.ErrEmptyName),
		errors.Is(err, domain_product.ErrInvalidProductID),
		errors.Is(err, domain_product.ErrInvalidAttribute):
		return app_error.New(app_error.KindValidation, errCodeProductInvalid, errMsgProductInvalid, err)
	case errors.Is(err, domain_product.ErrEmptySKUCode),
		errors.Is(err, domain_product.ErrDuplicateSKUCode),
		errors.Is(err, domain_product.ErrNoSKU),
		errors.Is(err, domain_product.ErrUnknownAttributeValue),
		errors.Is(err, domain_product.ErrInvalidSkuQuantity),
		errors.Is(err, domain_product.ErrNegativePrice),
		errors.Is(err, domain_product.ErrEmptyCurrency):
		return app_error.New(app_error.KindValidation, errCodeProductVariantInvalid, errMsgProductVariantInvalid, err)
	default:
		return err
	}
}
