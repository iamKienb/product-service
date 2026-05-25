package product

import "product-command-module/internal/domain/shared"

type Attribute struct {
	ID        shared.AttributeID
	ProductID shared.ProductID
	Name      string
}
