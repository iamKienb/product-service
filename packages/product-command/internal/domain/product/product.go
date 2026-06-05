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
	Variants        []ProductVariant

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

func (p *Product) AddAttribute(name string, values []string) error {
	if name == "" || len(values) == 0 {
		return ErrInvalidAttribute
	}
	for _, existing := range p.Attributes {
		if existing.Name == name {
			return ErrInvalidAttribute
		}
	}

	attributeID := shared.NewID[shared.AttributeID]()
	attribute := Attribute{
		ID:        attributeID,
		ProductID: p.ID,
		Name:      name,
	}

	p.Attributes = append(p.Attributes, attribute)

	for _, valName := range values {
		if valName == "" {
			return ErrInvalidAttribute
		}
		if _, exists := p.valueToID[valName]; exists {
			return ErrInvalidAttribute
		}

		valueID := shared.NewID[shared.AttributeValueID]()
		p.valueToID[valName] = valueID

		p.AttributeValues = append(p.AttributeValues, AttributeValue{
			ID:          valueID,
			AttributeID: attributeID,
			Name:        valName,
		})
	}

	return nil
}

func (p *Product) AddVariant(params ProductVariantParams, isDefault bool) (shared.SkuID, error) {
	if params.SkuCode == "" {
		return shared.SkuID{}, ErrEmptySKUCode
	}
	if params.Price <= 0 {
		return shared.SkuID{}, ErrNegativePrice
	}
	if params.Currency == "" {
		return shared.SkuID{}, ErrEmptyCurrency
	}
	for _, variant := range p.Variants {
		if variant.SkuCode == params.SkuCode {
			return shared.SkuID{}, ErrDuplicateSKUCode
		}
	}

	var valueIDs []shared.AttributeValueID

	for _, valueName := range params.AttributeValueNames {
		valueID, exists := p.valueToID[valueName]
		if !exists {
			return shared.SkuID{}, ErrUnknownAttributeValue
		}
		valueIDs = append(valueIDs, valueID)
	}

	skuID := shared.NewID[shared.SkuID]()
	variant := NewProductVariant(NewVariantParams{
		SkuID:             skuID,
		ProductID:         p.ID,
		ShopID:            p.ShopID,
		CreatedBy:         p.CreatedBy,
		CreatedAt:         p.CreatedAt,
		SkuCode:           params.SkuCode,
		Price:             params.Price,
		Currency:          params.Currency,
		ImageUrl:          params.ImageUrl,
		IsDefault:         isDefault,
		AttributeValueIDs: valueIDs,
	})

	p.Variants = append(p.Variants, *variant)
	p.updatePriceRange(variant.Price)

	return variant.SkuID, nil
}

func (p *Product) MarkAsCreated() {
	attrPayloads := make([]AttributePayload, 0, len(p.Attributes))
	for _, attr := range p.Attributes {
		attrValsPayloads := make([]AttributeValuePayload, 0)

		for _, val := range p.AttributeValues {
			if val.AttributeID == attr.ID {
				attrValsPayloads = append(attrValsPayloads, AttributeValuePayload{
					AttributeValueID: val.ID.String(),
					ValueName:        val.Name,
				})
			}
		}

		attrPayloads = append(attrPayloads, AttributePayload{
			AttributeID:   attr.ID.String(),
			AttributeName: attr.Name,
			Values:        attrValsPayloads,
		})
	}

	variantPayloads := make([]VariantPayload, 0, len(p.Variants))
	for _, variant := range p.Variants {
		attrValueStrings := make([]string, 0, len(variant.AttributeValueIDs))
		for _, id := range variant.AttributeValueIDs {
			attrValueStrings = append(attrValueStrings, id.String())
		}

		variantPayloads = append(variantPayloads, VariantPayload{
			SkuID:             variant.SkuID.String(),
			SkuCode:           variant.SkuCode,
			Price:             variant.Price,
			Currency:          variant.Currency,
			ImageUrl:          variant.ImageUrl,
			Status:            string(variant.Status),
			IsDefault:         variant.IsDefault,
			AttributeValueIDs: attrValueStrings,
		})

	}

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
		Attributes:  attrPayloads,
		Variants:    variantPayloads,
	})
}

func (p *Product) MarkAsDeleted() {
	p.AddEvent(ProductDeletedEvent{
		ProductID: p.ID,
		ShopID:    p.ShopID,
	})
}

func (p *Product) updatePriceRange(price int64) {
	if p.PriceMin == 0 || price < p.PriceMin {
		p.PriceMin = price
	}
	if price > p.PriceMax {
		p.PriceMax = price
	}
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
