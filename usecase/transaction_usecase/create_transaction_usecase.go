package transactionusecase

import (
	dto "be-shop-vision/dto/transaction"
	"context"
	"fmt"

	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (uc *TransactionUseCase) CreateTransaction(ctx context.Context, body dto.CreateTransactionDTO, storeID primitive.ObjectID) *entity.HttpError {
	var transactionProductAttr []model.TransactionProductAttribute
	totalPrice := 0

	for _, productDTO := range body.Data {
		productObjectID, _ := primitive.ObjectIDFromHex(productDTO.ProductID)

		product, err := uc.ProductService.FindOne(ctx, bson.M{"_id": productObjectID, "storeId": storeID})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return entity.BadRequest("Product tidak ditemukan")
			}
			return entity.InternalServerError(err.Error())
		}

		// Default pakai data produk utama
		itemName := product.Name
		itemPrice := 0
		itemStock := 0

		// Cek apakah pembelian ke variant (productDTO.VariantName ada isinya)
		if productDTO.VariantName != "" {
			foundVariant := false
			for _, variant := range product.Variants {
				if variant.Name == productDTO.VariantName {
					itemName = fmt.Sprintf("%s - %s", product.Name, variant.Name)
					itemPrice = variant.Price
					itemStock = variant.Stock
					foundVariant = true
					break
				}
			}

			if !foundVariant {
				return entity.BadRequest(fmt.Sprintf("Variant '%s' tidak ditemukan pada produk '%s'", productDTO.VariantName, product.Name))
			}
		} else {
			// Kalau bukan variant, pakai harga & stock produk utama
			if product.Price == nil || product.Stock == nil {
				return entity.BadRequest("Produk utama tidak memiliki harga atau stock")
			}
			itemPrice = *product.Price
			itemStock = *product.Stock
		}

		// Validasi stock
		if itemStock < productDTO.Quantity {
			return entity.BadRequest(fmt.Sprintf("Stock tidak mencukupi untuk %s", itemName))
		}

		// Tambahkan ke transaksi
		totalPriceProduct := itemPrice * productDTO.Quantity
		transactionProductAttr = append(transactionProductAttr, model.TransactionProductAttribute{
			ID:   product.ID,
			Name: itemName,
			Category: model.AttributeEmbedded{
				ID:   product.Category.ID,
				Name: product.Category.Name,
			},
			Price:      itemPrice,
			Quantity:   productDTO.Quantity,
			TotalPrice: totalPriceProduct,
		})

		totalPrice += totalPriceProduct
	}

	// Buat transaksi baru
	transaction := model.Transaction{
		ID:         primitive.NewObjectID(),
		Products:   transactionProductAttr,
		TotalPrice: totalPrice,
		StoreID:    storeID,
	}

	_, err := uc.TransactionService.Create(ctx, transaction)
	if err != nil {
		return entity.InternalServerError(err.Error())
	}

	// Setelah transaksi berhasil, update stock
	for _, productDTO := range body.Data {
		productObjectID, _ := primitive.ObjectIDFromHex(productDTO.ProductID)

		// Cari lagi untuk update stock
		product, err := uc.ProductService.FindOne(ctx, bson.M{"_id": productObjectID, "storeId": storeID})
		if err != nil {
			return entity.InternalServerError(err.Error())
		}

		// Update stock sesuai target (main product atau variant)
		if productDTO.VariantName != "" {
			// Update stock variant
			updateQuery := bson.M{
				"_id":      product.ID,
				"variants": bson.M{"$elemMatch": bson.M{"name": productDTO.VariantName}},
			}
			updateAction := bson.M{
				"$inc": bson.M{"variants.$.stock": -productDTO.Quantity},
			}

			_, err = uc.ProductService.UpdateOne(ctx, updateQuery, updateAction)
			if err != nil {
				return entity.InternalServerError(fmt.Sprintf("Gagal update stock untuk variant '%s': %s", productDTO.VariantName, err.Error()))
			}
		} else {
			// Update stock produk utama
			_, err = uc.ProductService.UpdateOne(ctx, bson.M{"_id": product.ID}, bson.M{"$inc": bson.M{"stock": -productDTO.Quantity}})
			if err != nil {
				return entity.InternalServerError(fmt.Sprintf("Gagal update stock untuk produk '%s': %s", product.Name, err.Error()))
			}
		}
	}

	return nil
}
