package events

const TopicProductDeleted = "product-service.product.deleted"

type ProductDeleted struct {
	ProductID string `json:"product_id"`
	ShopID    string `json:"shop_id"`
}
