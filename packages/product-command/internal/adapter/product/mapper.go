package product

import (
	"product-command-module/internal/application/commands/create_product"
	"product-command-module/internal/domain/shared"

	"github.com/iamKienb/api-contract/gen/product"
)

func ToCreateProductCommand(userID string, req *product.CreateProductsRequest) (create_product.Command, error) {
	parsedUserID, err := shared.ParseToRawID[shared.UserID](userID)
	if err != nil {
		return create_product.Command{}, shared.ErrInvalidUserID
	}

	shopID, err := shared.ParseToRawID[shared.ShopID](req.GetShopId())
	if err != nil {
		return create_product.Command{}, shared.ErrInvalidShopID
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
		ShopID:      shopID,
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
		Action:     req.GetAction(),
	}, nil
}

func ToCreateProductResponse(result *create_product.Result) *product.CreateProductsResponse {
	return &product.CreateProductsResponse{
		WorkflowId: result.ProductID,
	}
}
