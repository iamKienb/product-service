package product

import (
	"product-command-module/internal/application/services/product/i18n"
	domain_product "product-command-module/internal/domain/product"

	"github.com/iamKienb/go-core/app_error"
)

var productErrorMap = app_error.ServiceErrorMap{
	domain_product.ErrEmptyName:             {Kind: app_error.KindValidation, Msg: i18n.MsgProductNameEmpty},
	domain_product.ErrEmptySlug:             {Kind: app_error.KindValidation, Msg: i18n.MsgProductSlugEmpty},
	domain_product.ErrEmptySKUCode:          {Kind: app_error.KindValidation, Msg: i18n.MsgSkuCodeEmpty},
	domain_product.ErrEmptyCurrency:         {Kind: app_error.KindValidation, Msg: i18n.MsgProductCurrencyEmpty},
	domain_product.ErrInvalidProductID:      {Kind: app_error.KindValidation, Msg: i18n.MsgProductIdInvalid},
	domain_product.ErrInvalidAttribute:      {Kind: app_error.KindValidation, Msg: i18n.MsgProductAttributeInvalid},
	domain_product.ErrUnknownAttributeValue: {Kind: app_error.KindValidation, Msg: i18n.MsgProductAttributeValueUnknown},
	domain_product.ErrInvalidSkuQuantity:    {Kind: app_error.KindValidation, Msg: i18n.MsgSkuQuantityInvalid},
	domain_product.ErrNoSKU:                 {Kind: app_error.KindValidation, Msg: i18n.MsgProductNoSku},

	domain_product.ErrProductNotFound: {Kind: app_error.KindNotFound, Msg: i18n.MsgProductNotFound},

	domain_product.ErrProductSlugTaken: {Kind: app_error.KindConflict, Msg: i18n.MsgProductSlugTaken},
	domain_product.ErrDuplicateSKUCode: {Kind: app_error.KindConflict, Msg: i18n.MsgSkuCodeDuplicate},

	domain_product.ErrInvalidProductAction:    {Kind: app_error.KindForbidden, Msg: i18n.MsgProductActionForbidden},
	domain_product.ErrInvalidStatusTransition: {Kind: app_error.KindForbidden, Msg: i18n.MsgProductStatusTransitionInvalid},
}

func toApplicationError(err error) error {
	return app_error.WrapError(err, productErrorMap)
}
