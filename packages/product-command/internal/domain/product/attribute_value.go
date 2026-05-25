package product

import "product-command-module/internal/domain/shared"

type AttributeValue struct {
	ID          shared.AttributeValueID
	AttributeID shared.AttributeID
	Name        string
}
