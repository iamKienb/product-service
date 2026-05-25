package shared

import "github.com/google/uuid"

type CategoryID uuid.UUID
type ProductID uuid.UUID
type SkuID uuid.UUID
type AttributeID uuid.UUID
type AttributeValueID uuid.UUID

func NewID[T ~[16]byte]() T {
	return T(uuid.Must(uuid.NewV7()))
}

func (id CategoryID) String() string {
	return "cat_" + uuid.UUID(id).String()
}
func (id ProductID) String() string {
	return "prod_" + uuid.UUID(id).String()
}
func (id SkuID) String() string {
	return "sku_" + uuid.UUID(id).String()
}
func (id AttributeID) String() string {
	return "attr_" + uuid.UUID(id).String()
}
func (id AttributeValueID) String() string {
	return "val_" + uuid.UUID(id).String()
}

func (id CategoryID) RawID() uuid.UUID {
	return uuid.UUID(id)
}
func (id ProductID) RawID() uuid.UUID {
	return uuid.UUID(id)
}

func (id SkuID) RawID() uuid.UUID {
	return uuid.UUID(id)
}
func (id AttributeID) RawID() uuid.UUID {
	return uuid.UUID(id)
}
func (id AttributeValueID) RawID() uuid.UUID {
	return uuid.UUID(id)
}

type UserID uuid.UUID
type ShopID uuid.UUID

func (id UserID) String() string {
	return "user_" + uuid.UUID(id).String()
}

func (id ShopID) String() string {
	return "shop_" + uuid.UUID(id).String()
}

func (id UserID) RawID() uuid.UUID {
	return uuid.UUID(id)
}

func (id ShopID) RawID() uuid.UUID {
	return uuid.UUID(id)
}
