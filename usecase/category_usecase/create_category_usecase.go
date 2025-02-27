package categoryusecase

import (
	dto "be-shop-vision/dto/category"
	"context"

	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (uc *CategoryUseCase) CreateCategory(ctx context.Context, body dto.CreateCategoryDTO, userID primitive.ObjectID, storeID primitive.ObjectID) *entity.HttpError {
	_, err := uc.CategoryService.FindOne(ctx, bson.M{"name": body.Name, "storeId": storeID})
	if err != nil && err != mongo.ErrNoDocuments {
		return entity.InternalServerError(err.Error())
	}
	if err == nil {
		return entity.BadRequest("Nama kategori duplikat")
	}
	newCategory := model.Category{
		ID:      primitive.NewObjectID(),
		Name:    body.Name,
		StoreID: storeID,
	}

	_, errInsert := uc.CategoryService.Create(ctx, newCategory)

	if errInsert != nil {
		return entity.InternalServerError(errInsert.Error())
	}

	return nil
}
