package productusecase

import (
	dto "be-shop-vision/dto/product"
	"context"

	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (uc *ProductUseCase) CreateProduct(ctx context.Context, body dto.CreateProductDTO, storeID primitive.ObjectID, photos []string) *entity.HttpError {
	categoryObjectID, _ := primitive.ObjectIDFromHex(body.CategoryID)
	if len(photos) != 5 {
		return entity.BadRequest("Jumlah foto harus 5")
	}

	// Validasi coverPhoto dalam rentang index yang benar
	if body.CoverPhoto > len(photos) {
		return entity.BadRequest("Index coverPhoto tidak valid")
	}

	// Ambil foto berdasarkan coverPhoto index (karena coverPhoto dimulai dari 1, kurangi 1 untuk akses index array)
	coverPhoto := photos[body.CoverPhoto-1]

	category, err := uc.CategoryService.FindOne(ctx, bson.M{"_id": categoryObjectID, "storeId": storeID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.BadRequest("Kategori tidak ditemukan")
		}
		return entity.InternalServerError(err.Error())
	}

	product := model.Product{
		ID:         primitive.NewObjectID(),
		Name:       body.Name,
		StoreID:    storeID,
		Category:   model.AttributeEmbedded{ID: &category.ID, Name: &category.Name},
		Stock:      0,
		CoverPhoto: coverPhoto,
		Price:      body.Price,
	}

	_, err = uc.ProductService.Create(ctx, product)

	if err != nil {
		return entity.InternalServerError(err.Error())
	}
	key := 0
	for _, photo := range photos {
		key++
		productPhoto := model.ProductPhoto{Photo: photo, ProductID: product.ID, Key: key}
		_, err = uc.ProductPhotoService.Create(ctx, productPhoto)
	}

	return nil
}
