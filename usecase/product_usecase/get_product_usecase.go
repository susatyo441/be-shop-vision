package productusecase

import (
	"be-shop-vision/dto"
	productpipeline "be-shop-vision/pipeline/product"
	"context"

	utilDto "github.com/susatyo441/go-ta-utils/dto"
	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/functions"
	"github.com/susatyo441/go-ta-utils/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (uc *ProductUseCase) GetProductDetail(ctx context.Context, productID primitive.ObjectID, storeID primitive.ObjectID) (interface{}, *entity.HttpError) {
	// Ambil detail produk
	product, err := uc.ProductService.FindOne(ctx, bson.M{"_id": productID, "storeId": storeID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, entity.BadRequest("Produk tidak ditemukan")
		}
		return nil, entity.InternalServerError(err.Error())
	}

	// Ambil daftar foto produk
	productPhotos, err := uc.ProductPhotoService.Find(ctx, bson.M{"productId": productID})
	if err != nil {
		return nil, entity.InternalServerError("Gagal mengambil foto produk")
	}

	// Buat objek response langsung tanpa struct tambahan
	productDetail := map[string]interface{}{
		"_id":        product.ID,
		"name":       product.Name,
		"storeId":    product.StoreID,
		"category":   product.Category,
		"stock":      product.Stock,
		"coverPhoto": product.CoverPhoto,
		"photos":     productPhotos,
	}

	return productDetail, nil
}

func (u *ProductUseCase) GetProductList(ctx context.Context, query dto.PaginationQuery, storeID primitive.ObjectID) (*utilDto.PaginationResult[model.Product], *entity.HttpError) {

	var result []utilDto.PaginationResult[model.Product]

	// perform the aggregation query and handle any errors that occur during the process
	aggregateErr := u.ProductService.Aggregate(
		&result,
		ctx,
		productpipeline.GetProductsPipeline(query, storeID),
		options.Aggregate().SetCollation(&options.Collation{Strength: 3, Locale: "en"}),
	)

	if aggregateErr != nil {
		return nil, entity.InternalServerError(aggregateErr.Error())
	}

	response := functions.FormatPaginationResult(result)
	// return the first aggregation result
	return &response, nil
}
