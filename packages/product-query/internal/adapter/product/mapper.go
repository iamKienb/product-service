package product

import (
	"product-query-module/internal/application/queries/get_product_by_sku_ids"

	"github.com/iamKienb/api-contract/gen/product"
)

func toGetProductBySkuIDsQuery(req *product.GetProductsBySkuIDsRequest) (get_product_by_sku_ids.Query, error) {
	return get_product_by_sku_ids.Query{
		ShopID: req.GetShopId(),
		SkuIDs: req.GetSkuIds(),
	}, nil
}

func toGetProductBySkuIDsResponse(results []*get_product_by_sku_ids.Result) *product.GetProductsBySkuIDsResponse {
	items := make([]*product.SkuCheckoutDetail, 0, len(results))

	for _, result := range results {
		items = append(items, &product.SkuCheckoutDetail{
			SkuId:       result.SkuID,
			ProductId:   result.ProductID,
			ShopId:      result.ShopID,
			SkuCode:     result.SkuCode,
			ProductName: result.ProductName,
			Price:       result.Price,
			ImageUrl:    result.ImageURL,
			Status:      result.Status,
		})
	}

	return &product.GetProductsBySkuIDsResponse{
		Items: items,
	}
}
