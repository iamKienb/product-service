package product

import (
	"product-command-module/internal/domain/shared"
	"time"
)

type ProductStatus string

const (
	StatusDraft    ProductStatus = "DRAFT"
	StatusActive   ProductStatus = "ACTIVE"
	StatusARCHIVED ProductStatus = "ARCHIVED"
)

const Type = "PRODUCT"

type Product struct {
	ID          shared.ProductID
	ShopID      shared.ShopID
	Name        string
	Slug        string
	Description string
	Brand       string
	ThumbUrl    string
	VideoUrl    string

	PriceMin int64
	PriceMax int64

	Status     ProductStatus
	HasVariant bool

	CreatedBy shared.UserID
	UpdatedBy *shared.UserID

	CreatedAt time.Time
	UpdatedAt *time.Time

	Attributes      []Attribute
	AttributeValues []AttributeValue
	Variants        []Variant

	shared.EventEntity
	valueToID map[string]shared.AttributeValueID
}

func NewProduct(params NewProductParams) *Product {
	productID := shared.NewID[shared.ProductID]()
	now := time.Now().UTC()

	return &Product{
		ID:          productID,
		ShopID:      params.ShopID,
		Name:        params.Name,
		Slug:        params.Slug,
		Description: params.Description,
		Brand:       params.Brand,
		ThumbUrl:    params.ThumbUrl,
		VideoUrl:    params.VideoUrl,
		Status:      StatusDraft,
		HasVariant:  params.HasVariant,
		CreatedBy:   params.UserID,
		CreatedAt:   now,
		valueToID:   make(map[string]shared.AttributeValueID),
	}

}

// func (a *Product) UpdateProductInfo(params UpdateProductParam) {
// 	a.Name = params.Name
// 	a.Description = params.Description
// 	a.ThumbUrl = params.ThumbUrl
// 	a.VideoUrl = params.VideoUrl

// 	now := time.Now().UTC()
// 	a.UpdatedAt = &now

// 		a.AddEvent(ProductUpdatedEvent{
// 			ProductID:   a.ID,
// 			Name:        a.Name,
// 			Description: a.Description,
// 			UpdatedAt:   *a.UpdatedAt,
// 		})
// 	}

func (p *Product) AddAttribute(name string, values []string) {
	attributeID := shared.NewID[shared.AttributeID]()
	attribute := Attribute{
		ID:        attributeID,
		ProductID: p.ID,
		Name:      name,
	}

	p.AddEvent(AttributeCreatedEvent{
		AttributeID:   attributeID,
		ProductID:     p.ID,
		AttributeName: attribute.Name,
	})

	p.Attributes = append(p.Attributes, attribute)

	for _, valName := range values {
		valueID := shared.NewID[shared.AttributeValueID]()
		p.valueToID[valName] = valueID

		attributeValue := AttributeValue{
			ID:          valueID,
			AttributeID: attributeID,
			Name:        valName,
		}

		p.AddEvent(AttributeValueCreatedEvent{
			AttributeValueID: valueID,
			AttributeID:      attributeID,
			ValueName:        attributeValue.Name,
		})
		p.AttributeValues = append(p.AttributeValues, attributeValue)
	}
}

func (p *Product) AddVariant(params VariantParam, isDefault bool) shared.SkuID {
	skuID := shared.NewID[shared.SkuID]()
	var valueIDs []shared.AttributeValueID

	for _, valueName := range params.AttributeValueNames {
		if valueID, exists := p.valueToID[valueName]; exists {
			valueIDs = append(valueIDs, valueID)
		}
	}

	variant := Variant{
		SkuID:             skuID,
		ProductID:         p.ID,
		ShopID:            p.ShopID,
		SkuCode:           params.SkuCode,
		Price:             params.Price,
		Currency:          params.Currency,
		ImageUrl:          params.ImageUrl,
		IsDefault:         isDefault,
		AttributeValueIDs: valueIDs,

		CreatedBy: p.CreatedBy,
		CreatedAt: p.CreatedAt,
	}

	p.AddEvent(VariantCreatedEvent{
		SkuID:             skuID,
		ProductID:         variant.ProductID,
		ShopID:            variant.ShopID,
		SkuCode:           variant.SkuCode,
		Price:             variant.Price,
		Currency:          variant.Currency,
		ImageUrl:          variant.ImageUrl,
		AttributeValueIDs: variant.AttributeValueIDs,

		CreatedBy: variant.CreatedBy,
		CreatedAt: variant.CreatedAt,
	})

	p.Variants = append(p.Variants, variant)
	return skuID
}

func (p *Product) MarkAsCreated() {
	p.AddEvent(ProductCreatedEvent{
		ProductID:   p.ID,
		ShopID:      p.ShopID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Brand:       p.Brand,
		ThumbUrl:    p.ThumbUrl,
		VideoUrl:    p.VideoUrl,
		PriceMin:    p.PriceMin,
		PriceMax:    p.PriceMax,
		Status:      p.Status,
		HasVariant:  p.HasVariant,
		CreatedBy:   p.CreatedBy,
		CreatedAt:   p.CreatedAt,
	})
}

func (p *Product) MarkAsDeleted() {
	p.AddEvent(ProductDeletedEvent{
		ProductID: p.ID,
		ShopID:    p.ShopID,
	})
}

func (p *Product) recalculatePriceRange() {
	if len(p.Variants) == 0 {
		return
	}

	min := p.Variants[0].Price
	max := p.Variants[0].Price

	for _, variant := range p.Variants {
		if variant.Price < min {
			min = variant.Price
		}
		if variant.Price > max {
			max = variant.Price
		}
	}

	p.PriceMin = min
	p.PriceMax = max
}

// func (p *Product) UpdateVariant(params updateVariantParams) error {
// 	if p.Status == "BANNED" {
// 		return errors.New("product is banned, cannot update price")
// 	}

// 	found := false
// 	for i, v := range p.Variants {
// 		if v.SkuID == skuID {
// 			p.Variants[i].Price = params.Price
// 			p.Variants[i].ImageUrl = params.ImageUrl
// 			found = true
// 			break
// 		}
// 	}
// 	if !found {
// 		return errors.New("variant not found in this product")
// 	}

// 	// 3. Sinh Event riêng lẻ
// 	p.AddEvent(VariantUpdatedEvent{
// 		ProductID: p.ID,
// 		SkuID:     skuID,
// 		Price:     Price,
// 		ImageUrl:  ImageUrl,
// 	})
// 	return nil
// }

func (p *Product) FlushEvents() []shared.DomainEvent {
	var domainEvents []shared.DomainEvent

	domainEvents = append(domainEvents, p.Flush()...)
	p.ClearEvent()

	return domainEvents
}

func (p Product) Type() string {
	return "PRODUCT"
}
