package product

import (
	"product-command-module/internal/application/services/product/i18n"
	"product-command-module/internal/domain/product"

	"github.com/iamKienb/go-core/app_error"
)

var productErrorMap = app_error.ServiceErrorMap{
	product.ErrProductSlugTaken: {Kind: app_error.KindValidation, Msg: i18n.MsgSlugTaken},
}

func (s *productService) wrapError(err error) error {
	return app_error.WrapError(err, productErrorMap)
}
