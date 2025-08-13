package productusecase

import (
	"be-shop-vision/dto"
	productdto "be-shop-vision/dto/product"
	"context"
	"mime/multipart"

	"github.com/susatyo441/go-ta-utils/db"
	utilDto "github.com/susatyo441/go-ta-utils/dto"
	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/model"
	utilservice "github.com/susatyo441/go-ta-utils/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IProductUseCase interface {
	CreateProduct(ctx context.Context, body productdto.CreateProductDTO, storeID primitive.ObjectID, files map[string]*multipart.FileHeader) *entity.HttpError
	BulkDeleteProducts(ctx context.Context, body dto.ArrayOfIdDTO, storeID primitive.ObjectID) *entity.HttpError
	UpdateProduct(ctx context.Context, productID primitive.ObjectID, body productdto.UpdateProductDTO, storeID primitive.ObjectID, files map[string]*multipart.FileHeader) *entity.HttpError
	GetProductDetail(ctx context.Context, productID primitive.ObjectID, storeID primitive.ObjectID) (interface{}, *entity.HttpError)
	GetProductList(ctx context.Context, query dto.PaginationQuery, storeID primitive.ObjectID) (*utilDto.PaginationResult[model.Product], *entity.HttpError)
	ExportAllData(ctx context.Context) ([]model.Product, []model.Category, []model.ProductPhoto, []model.Transaction, []model.Store, []model.User, error)
	UpdateProductStock(ctx context.Context, bodies []productdto.UpdateProductStockDTO, storeID primitive.ObjectID) *entity.HttpError
}

type ProductUseCase struct {
	ProductService      utilservice.Service[model.Product]
	CategoryService     utilservice.Service[model.Category]
	ProductPhotoService utilservice.Service[model.ProductPhoto]
	TransactionService  utilservice.Service[model.Transaction]
	UserService         utilservice.Service[model.User]
	StoreService        utilservice.Service[model.Store]
}

func MakeProductUseCase() IProductUseCase {
	return &ProductUseCase{
		CategoryService:     utilservice.ShopVisionService[model.Category](db.CategoryModelName),
		ProductService:      utilservice.ShopVisionService[model.Product](db.ProductModelName),
		ProductPhotoService: utilservice.ShopVisionService[model.ProductPhoto](db.ProductPhotoModelName),
		TransactionService:  utilservice.ShopVisionService[model.Transaction](db.TransactionsModelName),
		UserService:         utilservice.ShopVisionService[model.User](db.UserModelName),
		StoreService:        utilservice.ShopVisionService[model.Store](db.StoreModelName),
	}
}
