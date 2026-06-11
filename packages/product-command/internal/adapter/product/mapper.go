package product

import (
	"product-command-module/internal/application/commands/create_product"
	"product-command-module/internal/domain/shared"

	"github.com/iamKienb/api-contract/gen/product"
	"github.com/iamKienb/go-core/app_error"
)

func ToCreateProductCommand(userID string, req *product.CreateProductsRequest) (create_product.Command, error) {
	parsedUserID, err := parseUserID(userID)
	if err != nil {
		return create_product.Command{}, err
	}

	parsedShopID, err := parseShopID(req.GetShopId())
	if err != nil {
		return create_product.Command{}, err
	}

	attributes := make([]create_product.ProductAttribute, 0, len(req.GetAttributes()))
	if req.GetHasVariant() {
		for _, attributeReq := range req.GetAttributes() {
			attributes = append(attributes, create_product.ProductAttribute{
				Name:   attributeReq.GetName(),
				Values: attributeReq.GetValues(),
			})

		}
	}

	variants := make([]create_product.ProductVariant, 0, len(req.GetVariants()))
	for _, variantReq := range req.GetVariants() {
		variants = append(variants, create_product.ProductVariant{
			SkuCode:             variantReq.GetSkuCode(),
			Price:               variantReq.GetPrice(),
			Currency:            variantReq.GetCurrency(),
			ImageUrl:            variantReq.GetImageUrl(),
			AttributeValueNames: variantReq.GetAttributeValueNames(),
			Quantity:            variantReq.GetQuantity(),
		})

	}

	return create_product.Command{
		ShopID:      parsedShopID,
		UserID:      parsedUserID,
		Name:        req.GetName(),
		Slug:        req.GetSlug(),
		Description: req.GetDescription(),
		Brand:       req.GetBrand(),

		ThumbUrl: req.GetThumbUrl(),
		VideoUrl: req.GetVideoUrl(),

		Attributes: attributes,
		Variants:   variants,
		HasVariant: req.HasVariant,
		Status:     req.GetStatus(),
	}, nil
}

func ToCreateProductResponse(result *create_product.Result) *product.CreateProductsResponse {
	return &product.CreateProductsResponse{
		WorkflowId: result.ProductID.String(),
	}
}

func parseUserID(value string) (shared.UserID, error) {
	parsed, err := shared.ParseToRawID[shared.UserID](value)
	if err != nil {
		return parsed, app_error.New(app_error.KindValidation, "user_invalid", "invalid user id", err)
	}

	return parsed, nil
}

func parseShopID(value string) (shared.ShopID, error) {
	parsed, err := shared.ParseToRawID[shared.ShopID](value)
	if err != nil {
		return parsed, app_error.New(app_error.KindValidation, "shop_invalid", "invalid shop id", err)
	}

	return parsed, nil
}

func parseSkuID(value string) (shared.SkuID, error) {
	parsed, err := shared.ParseToRawID[shared.SkuID](value)
	if err != nil {
		return parsed, app_error.New(app_error.KindValidation, "sku_invalid", "invalid sku id", err)
	}

	return parsed, nil
}
