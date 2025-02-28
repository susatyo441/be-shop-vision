package categoryusecase

import (
	dto "be-shop-vision/dto/category"
	"context"

	"github.com/susatyo441/go-ta-utils/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (uc *CategoryUseCase) UpdateCategory(ctx context.Context, body dto.CreateCategoryDTO, categoryId primitive.ObjectID, storeID primitive.ObjectID) *entity.HttpError {
	_, err := uc.CategoryService.FindOne(ctx, bson.M{"name": body.Name, "storeId": storeID})

	if err != nil && err != mongo.ErrNoDocuments {
		return entity.InternalServerError(err.Error())
	}

	if err == nil {
		return entity.BadRequest("Nama kategori duplikat")
	}

	_, err = uc.CategoryService.FindOneAndUpdate(ctx, bson.M{"_id": categoryId}, bson.M{"$set": bson.M{
		"name": body.Name,
	}})

	if err != nil {
		if err != mongo.ErrNoDocuments {
			return entity.BadRequest("tidak ditemukan category")
		}
		return entity.InternalServerError(err.Error())
	}

	return nil
}
