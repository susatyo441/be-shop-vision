package categoryusecase

import (
	"be-shop-vision/dto"
	categorydto "be-shop-vision/dto/category"
	"context"

	"github.com/susatyo441/go-ta-utils/db"
	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/model"
	utilservice "github.com/susatyo441/go-ta-utils/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICategoryUseCase interface {
	CreateCategory(ctx context.Context, body categorydto.CreateCategoryDTO, storeID primitive.ObjectID) *entity.HttpError
	UpdateCategory(ctx context.Context, body categorydto.CreateCategoryDTO, categoryId primitive.ObjectID, storeID primitive.ObjectID) *entity.HttpError
	GetCategoryOptions(ctx context.Context, storeID primitive.ObjectID, query dto.PaginationQuery) (*categorydto.GetCategoryOptionsResponse, *entity.HttpError)
	BulkDeleteCategories(ctx context.Context, body dto.ArrayOfIdDTO, storeID primitive.ObjectID) *entity.HttpError
}

type CategoryUseCase struct {
	CategoryService utilservice.Service[model.Category]
}

func MakeCategoryUseCase() ICategoryUseCase {
	return &CategoryUseCase{

		CategoryService: utilservice.ShopVisionService[model.Category](db.CategoryModelName),
	}
}
