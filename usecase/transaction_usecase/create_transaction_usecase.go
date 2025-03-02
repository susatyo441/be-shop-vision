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
	productObjectID, _ := primitive.ObjectIDFromHex(body.ProductID)

	product, err := uc.ProductService.FindOne(ctx, bson.M{"_id": productObjectID, "storeId": storeID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.BadRequest("Product tidak ditemukan")
		}
		return entity.InternalServerError(err.Error())
	}

	if product.Stock < body.Quantity {
		return entity.BadRequest("Stock tidak mencukupi")
	}

	transaction := model.Transaction{
		ID:         primitive.NewObjectID(),
		Product:    model.AttributeEmbedded{ID: &product.ID, Name: &product.Name},
		Category:   model.AttributeEmbedded{ID: product.Category.ID, Name: product.Category.Name},
		Price:      product.Price,
		Quantity:   body.Quantity,
		TotalPrice: product.Price * body.Quantity,
		StoreID:    storeID,
	}

	_, err = uc.TransactionService.Create(ctx, transaction)

	if err != nil {
		return entity.InternalServerError(err.Error())
	}

	_, err = uc.ProductService.UpdateOne(ctx, bson.M{"_id": product.ID}, bson.M{"$inc": bson.M{"stock": -body.Quantity}})
	if err != nil {
		return entity.InternalServerError(err.Error())
	}

	return nil
}
