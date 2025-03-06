package transactionusecase

import (
	dto "be-shop-vision/dto/transaction"
	"context"

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

		if product.Stock < productDTO.Quantity {
			return entity.BadRequest("Stock tidak mencukupi")
		}
		totalPriceProduct := product.Price * productDTO.Quantity
		transactionProductAttr = append(transactionProductAttr, model.TransactionProductAttribute{
			ID:   product.ID,
			Name: product.Name,
			Category: model.AttributeEmbedded{
				ID:   product.Category.ID,
				Name: product.Category.Name,
			},
			Price:      product.Price,
			Quantity:   productDTO.Quantity,
			TotalPrice: totalPriceProduct,
		})
		totalPrice += totalPriceProduct
	}

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

	for _, productDTO := range body.Data {
		productObjectID, _ := primitive.ObjectIDFromHex(productDTO.ProductID)
		_, err = uc.ProductService.UpdateOne(ctx, bson.M{"_id": productObjectID}, bson.M{"$inc": bson.M{"stock": -productDTO.Quantity}})
		if err != nil {
			return entity.InternalServerError(err.Error())
		}
	}

	return nil
}
