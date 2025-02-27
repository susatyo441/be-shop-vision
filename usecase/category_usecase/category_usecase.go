package categoryusecase

import (
	dto "be-shop-vision/dto/category"
	"context"

	"github.com/susatyo441/go-ta-utils/db"
	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/model"
	utilservice "github.com/susatyo441/go-ta-utils/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICategoryUseCase interface {
	CreateCategory(ctx context.Context, body dto.CreateCategoryDTO, userID primitive.ObjectID, storeID primitive.ObjectID) *entity.HttpError
}

type CategoryUseCase struct {
	CategoryService utilservice.Service[model.Category]
}

func MakeCategoryUseCase() ICategoryUseCase {
	return &CategoryUseCase{

		CategoryService: utilservice.ShopVisionService[model.Category](db.CategoryModelName),
	}
}
