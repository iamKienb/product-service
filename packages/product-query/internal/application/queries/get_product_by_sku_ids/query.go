package get_product_by_sku_ids

import "context"

type Query struct {
	ShopID string
	SkuIDs []string
}

type Result struct {
	SkuID       string
	ProductID   string
	ShopID      string
	SkuCode     string
	ProductName string
	Price       int64
	ImageURL    string
	Status      string
}

type Executor interface {
	Execute(ctx context.Context, qry Query) ([]*Result, error)
}
