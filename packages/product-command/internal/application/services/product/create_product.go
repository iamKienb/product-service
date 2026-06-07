package product

import (
	"context"
	"errors"
	"sort"
	"strings"

	"product-command-module/internal/application/commands/create_product"
	"product-command-module/internal/application/port"
	"product-command-module/internal/domain/product"
)

const (
	createProductActionDraft   = "DRAFT"
	createProductActionPublish = "PUBLISH"
)

func (s *productService) CreateProduct(ctx context.Context, cmd create_product.Command) (*create_product.Result, error) {
	if err := validateCreateProductCommand(cmd); err != nil {
		return nil, err
	}

	existing, err := s.findExistingProductBySlug(ctx, cmd)
	if err == nil {
		return createProductResultFromExisting(cmd, existing)
	}
	if !errors.Is(err, product.ErrProductNotFound) {
		return nil, err
	}

	newProduct := product.NewProduct(product.NewProductParams{
		ShopID:      cmd.ShopID,
		UserID:      cmd.UserID,
		Name:        cmd.Name,
		Slug:        cmd.Slug,
		Description: cmd.Description,
		Brand:       cmd.Brand,
		ThumbUrl:    cmd.ThumbUrl,
		VideoUrl:    cmd.VideoUrl,
		HasVariant:  cmd.HasVariant,
	})

	if cmd.HasVariant {
		for _, attr := range cmd.Attributes {
			if err := newProduct.AddAttribute(attr.Name, attr.Values); err != nil {
				return nil, err
			}
		}
	}

	skuItems := make([]create_product.SkuItem, 0, len(cmd.Variants))
	for index, variant := range cmd.Variants {
		isDefault := index == 0

		variantParam := product.ProductVariantParams{
			SkuCode:             variant.SkuCode,
			Price:               variant.Price,
			Currency:            variant.Currency,
			ImageUrl:            variant.ImageUrl,
			AttributeValueNames: variant.AttributeValueNames,
		}

		skuID, err := newProduct.AddVariant(variantParam, isDefault)
		if err != nil {
			return nil, err
		}
		skuItems = append(skuItems, create_product.SkuItem{
			SkuID:    skuID,
			Quantity: variant.Quantity,
		})
	}

	action, err := normalizeCreateProductAction(cmd.Action)
	if err != nil {
		return nil, err
	}
	if action == createProductActionPublish {
		if err := newProduct.Publish(); err != nil {
			return nil, err
		}
	}

	newProduct.MarkAsCreated()

	if err := s.txManager.WithTx(ctx, func(ctx context.Context) error {
		if err := s.productRepo.CreateProduct(ctx, newProduct); err != nil {
			return err
		}

		if productEvents := newProduct.FlushEvents(); len(productEvents) > 0 {
			if err := s.outboxService.Publish(ctx, port.OutboxParam{
				AggregateID:   newProduct.ID.RawID(),
				AggregateType: newProduct.Type(),
				Events:        productEvents,
			}); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		if errors.Is(err, product.ErrProductSlugTaken) {
			existing, readErr := s.productRepo.GetProductByShopAndSlug(ctx, cmd.ShopID, cmd.Slug)
			if readErr == nil {
				return createProductResultFromExisting(cmd, existing)
			}
		}
		return nil, err
	}

	if s.productCache != nil {
		bgCtx := context.WithoutCancel(ctx)
		go func() {
			_ = s.productCache.RememberSlug(bgCtx, cmd.Slug)
		}()
	}

	return &create_product.Result{
		ShopID:    cmd.ShopID,
		ProductID: newProduct.ID,
		SkuItems:  skuItems,
	}, nil
}

func validateCreateProductCommand(cmd create_product.Command) error {
	if strings.TrimSpace(cmd.Name) == "" {
		return product.ErrEmptyName
	}
	if strings.TrimSpace(cmd.Slug) == "" {
		return product.ErrEmptySlug
	}
	if len(cmd.Variants) == 0 {
		return product.ErrNoSKU
	}
	if _, err := normalizeCreateProductAction(cmd.Action); err != nil {
		return err
	}

	seenAttributes := make(map[string]struct{}, len(cmd.Attributes))
	seenAttributeValues := make(map[string]struct{})
	for _, attr := range cmd.Attributes {
		if strings.TrimSpace(attr.Name) == "" || len(attr.Values) == 0 {
			return product.ErrInvalidAttribute
		}
		if _, exists := seenAttributes[attr.Name]; exists {
			return product.ErrInvalidAttribute
		}
		seenAttributes[attr.Name] = struct{}{}
		for _, value := range attr.Values {
			if strings.TrimSpace(value) == "" {
				return product.ErrInvalidAttribute
			}
			if _, exists := seenAttributeValues[value]; exists {
				return product.ErrInvalidAttribute
			}
			seenAttributeValues[value] = struct{}{}
		}
	}

	seenSkuCodes := make(map[string]struct{}, len(cmd.Variants))
	for _, variant := range cmd.Variants {
		if variant.Quantity < 0 {
			return product.ErrInvalidSkuQuantity
		}
		if _, exists := seenSkuCodes[variant.SkuCode]; exists {
			return product.ErrDuplicateSKUCode
		}
		seenSkuCodes[variant.SkuCode] = struct{}{}
	}
	return nil
}

func normalizeCreateProductAction(action string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(action))
	switch normalized {
	case "", createProductActionPublish:
		return createProductActionPublish, nil
	case createProductActionDraft:
		return createProductActionDraft, nil
	default:
		return "", product.ErrInvalidProductAction
	}
}

func statusForCreateProductAction(action string) (product.ProductStatus, error) {
	normalized, err := normalizeCreateProductAction(action)
	if err != nil {
		return "", err
	}
	if normalized == createProductActionDraft {
		return product.StatusDraft, nil
	}
	return product.StatusActive, nil
}

func (s *productService) findExistingProductBySlug(ctx context.Context, cmd create_product.Command) (*product.Product, error) {
	isDuplicateSlug, err := s.productRepo.CheckSlugExists(ctx, cmd.ShopID, cmd.Slug)
	if err != nil {
		return nil, err
	}
	if !isDuplicateSlug {
		return nil, product.ErrProductNotFound
	}

	return s.productRepo.GetProductByShopAndSlug(ctx, cmd.ShopID, cmd.Slug)
}

func createProductResultFromExisting(cmd create_product.Command, existing *product.Product) (*create_product.Result, error) {
	if !matchesCreateProductCommand(cmd, existing) {
		return nil, product.ErrProductSlugTaken
	}

	variantsByCode := make(map[string]product.ProductVariant, len(existing.Variants))
	for _, variant := range existing.Variants {
		variantsByCode[variant.SkuCode] = variant
	}

	skuItems := make([]create_product.SkuItem, 0, len(cmd.Variants))
	for _, requested := range cmd.Variants {
		variant, exists := variantsByCode[requested.SkuCode]
		if !exists {
			return nil, product.ErrProductSlugTaken
		}
		skuItems = append(skuItems, create_product.SkuItem{
			SkuID:    variant.SkuID,
			Quantity: requested.Quantity,
		})
	}

	return &create_product.Result{
		ShopID:    existing.ShopID,
		ProductID: existing.ID,
		SkuItems:  skuItems,
	}, nil
}

func matchesCreateProductCommand(cmd create_product.Command, existing *product.Product) bool {
	expectedStatus, err := statusForCreateProductAction(cmd.Action)
	if err != nil {
		return false
	}

	if existing.ShopID != cmd.ShopID ||
		existing.Name != cmd.Name ||
		existing.Slug != cmd.Slug ||
		existing.Description != cmd.Description ||
		existing.Brand != cmd.Brand ||
		existing.ThumbUrl != cmd.ThumbUrl ||
		existing.VideoUrl != cmd.VideoUrl ||
		existing.Status != expectedStatus ||
		existing.HasVariant != cmd.HasVariant {
		return false
	}

	if !matchesAttributes(cmd.Attributes, existing) {
		return false
	}

	return matchesVariants(cmd.Variants, existing)
}

func matchesAttributes(requested []create_product.ProductAttribute, existing *product.Product) bool {
	if len(requested) != len(existing.Attributes) {
		return false
	}

	valuesByAttributeID := make(map[productAttributeID][]string, len(existing.Attributes))
	for _, value := range existing.AttributeValues {
		attributeID := productAttributeID(value.AttributeID.String())
		valuesByAttributeID[attributeID] = append(valuesByAttributeID[attributeID], value.Name)
	}

	existingValuesByName := make(map[string][]string, len(existing.Attributes))
	for _, attribute := range existing.Attributes {
		attributeID := productAttributeID(attribute.ID.String())
		existingValuesByName[attribute.Name] = sortedStrings(valuesByAttributeID[attributeID])
	}

	seenRequestedAttributes := make(map[string]struct{}, len(requested))
	for _, attr := range requested {
		if _, exists := seenRequestedAttributes[attr.Name]; exists {
			return false
		}
		seenRequestedAttributes[attr.Name] = struct{}{}

		existingValues, exists := existingValuesByName[attr.Name]
		if !exists {
			return false
		}
		if !sameStrings(existingValues, sortedStrings(attr.Values)) {
			return false
		}
	}

	return true
}

func matchesVariants(requested []create_product.ProductVariant, existing *product.Product) bool {
	if len(requested) != len(existing.Variants) {
		return false
	}

	attributeValueNameByID := make(map[string]string, len(existing.AttributeValues))
	for _, value := range existing.AttributeValues {
		attributeValueNameByID[value.ID.String()] = value.Name
	}

	existingBySkuCode := make(map[string]product.ProductVariant, len(existing.Variants))
	for _, variant := range existing.Variants {
		existingBySkuCode[variant.SkuCode] = variant
	}

	seenRequestedSkuCodes := make(map[string]struct{}, len(requested))
	for _, req := range requested {
		if _, exists := seenRequestedSkuCodes[req.SkuCode]; exists {
			return false
		}
		seenRequestedSkuCodes[req.SkuCode] = struct{}{}

		variant, exists := existingBySkuCode[req.SkuCode]
		if !exists ||
			variant.Price != req.Price ||
			variant.Currency != req.Currency ||
			variant.ImageUrl != req.ImageUrl {
			return false
		}

		existingValueNames := make([]string, 0, len(variant.AttributeValueIDs))
		for _, id := range variant.AttributeValueIDs {
			existingValueNames = append(existingValueNames, attributeValueNameByID[id.String()])
		}
		if !sameStrings(sortedStrings(req.AttributeValueNames), sortedStrings(existingValueNames)) {
			return false
		}
	}

	return true
}

func sortedStrings(values []string) []string {
	result := append([]string(nil), values...)
	sort.Strings(result)
	return result
}

func sameStrings(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

type productAttributeID string
