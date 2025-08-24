package productusecase

import (
	"context"

	dto "be-shop-vision/dto/product"

	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/functions"
	"github.com/susatyo441/go-ta-utils/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (uc *ProductUseCase) UpdateProductStock(ctx context.Context, bodies []dto.UpdateProductStockDTO, storeID primitive.ObjectID) *entity.HttpError {
	for _, body := range bodies {
		productID, errObjectID := primitive.ObjectIDFromHex(body.ProductID)
		if errObjectID != nil {
			return entity.BadRequest("ID produk tidak valid")
		}
		// 1️⃣ Cari produk lama
		product, err := uc.ProductService.FindOne(ctx, bson.M{"_id": productID, "storeId": storeID})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return entity.BadRequest("Produk tidak ditemukan")
			}
			return entity.InternalServerError(err.Error())
		}

		variants := make([]model.ProductVariantsAttr, len(product.Variants))
		for i, variantDTO := range product.Variants {
			variants[i] = model.ProductVariantsAttr{
				Name:         variantDTO.Name,
				Price:        variantDTO.Price,
				CapitalPrice: nil,
				Stock:        functions.Ternary(*body.Variant == variantDTO.Name, body.Stock, variantDTO.Stock),
			}
		}

		if body.Variant == nil || *body.Variant == "" {
			product.Stock = functions.MakePointer(body.Stock)

		}
		// 7️⃣ Update produk ke database
		updateData := bson.M{
			"stock":    product.Stock,
			"variants": variants,
		}

		if _, err := uc.ProductService.UpdateOne(ctx, bson.M{"_id": productID, "storeId": storeID}, bson.M{"$set": updateData}); err != nil {
			return entity.InternalServerError("Gagal update produk: " + err.Error())
		}
	}

	return nil
}
