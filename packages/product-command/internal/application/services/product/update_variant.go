package product

// func (s *productService) UpdateVariant(ctx context.Context) error {

// 		// 1. REPO TỐI ƯU: Hàm này CHỈ SELECT bảng Product và bảng Variant
// 		// Tuyệt đối KHÔNG SELECT bảng Attribute, giúp giảm tải DB tối đa
// 		product, err := s.productRepo.FindByIDWithVariants(ctx, productID)
// 		if err != nil {
// 			return err
// 		}

// 		if err := productAgg.UpdateVariant(updateVariantParams); err != nil {
// 			return err
// 		}

// 		// 3. REPO TỐI ƯU: Chỉ chạy đúng 1 câu lệnh UPDATE bảng variants
// 		// Không hề động chạm hay ghi đè lại bảng products hay attributes
// 		if err := s.productRepo.UpdateProductVariant(ctx, tx, skuID, newPrice); err != nil {
// 			return err
// 		}

// 		// 4. Publish Event
// 		return s.publisher.Publish(ctx, productAgg.PullEvents())
// }
