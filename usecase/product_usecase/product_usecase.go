package productusecase

import (
	"be-shop-vision/dto"
	productdto "be-shop-vision/dto/product"
	"context"

	"github.com/susatyo441/go-ta-utils/db"
	utilDto "github.com/susatyo441/go-ta-utils/dto"
	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/model"
	utilservice "github.com/susatyo441/go-ta-utils/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IProductUseCase interface {
	CreateProduct(ctx context.Context, body productdto.CreateProductDTO, storeID primitive.ObjectID, photos []string) *entity.HttpError
	BulkDeleteProducts(ctx context.Context, body dto.ArrayOfIdDTO, storeID primitive.ObjectID) *entity.HttpError
	UpdateProduct(ctx context.Context, productID primitive.ObjectID, body productdto.UpdateProductDTO, storeID primitive.ObjectID) *entity.HttpError
	GetProductDetail(ctx context.Context, productID primitive.ObjectID, storeID primitive.ObjectID) (interface{}, *entity.HttpError)
	GetProductList(ctx context.Context, query dto.PaginationQuery, storeID primitive.ObjectID) (*utilDto.PaginationResult[model.Product], *entity.HttpError)
}

type ProductUseCase struct {
	ProductService      utilservice.Service[model.Product]
	CategoryService     utilservice.Service[model.Category]
	ProductPhotoService utilservice.Service[model.ProductPhoto]
}

func MakeProductUseCase() IProductUseCase {
	return &ProductUseCase{
		CategoryService:     utilservice.ShopVisionService[model.Category](db.CategoryModelName),
		ProductService:      utilservice.ShopVisionService[model.Product](db.ProductModelName),
		ProductPhotoService: utilservice.ShopVisionService[model.ProductPhoto](db.ProductPhotoModelName),
	}
}
