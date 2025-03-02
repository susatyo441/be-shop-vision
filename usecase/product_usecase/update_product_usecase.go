package productusecase

import (
	dto "be-shop-vision/dto/product"
	"context"
	"fmt"

	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/functions"
	"github.com/susatyo441/go-ta-utils/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (uc *ProductUseCase) UpdateProduct(ctx context.Context, productID primitive.ObjectID, body dto.UpdateProductDTO, storeID primitive.ObjectID) *entity.HttpError {
	// Cari produk yang akan diperbarui
	_, err := uc.ProductService.FindOne(ctx, bson.M{"_id": productID, "storeId": storeID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.BadRequest("Produk tidak ditemukan")
		}
		return entity.InternalServerError(err.Error())
	}

	// Ambil kategori untuk validasi
	categoryObjectID, _ := primitive.ObjectIDFromHex(body.CategoryID)
	category, err := uc.CategoryService.FindOne(ctx, bson.M{"_id": categoryObjectID, "storeId": storeID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.BadRequest("Kategori tidak ditemukan")
		}
		return entity.InternalServerError(err.Error())
	}

	// Validasi CoverPhoto index (1-based index)
	photos := []string{
		body.Photos.FirstImage,
		body.Photos.SecondImage,
		body.Photos.ThirdImage,
		body.Photos.FourthImage,
		body.Photos.FifthImage,
	}

	key := 0
	for _, photo := range photos {
		key++
		if photo != "" {
			productPhoto, errUpdate := uc.ProductPhotoService.FindOneAndUpdate(ctx, bson.M{"productId": productID, "key": key}, bson.M{"$set": bson.M{"photo": photo}}, options.FindOneAndUpdate().SetReturnDocument(options.Before))
			if errUpdate != nil {
				return entity.InternalServerError(errUpdate.Error())
			}

			err = functions.DeleteImage(fmt.Sprintf("/%s", storeID.Hex()), productPhoto.Photo)
			if err != nil {
				return entity.InternalServerError(err.Error())
			}

		}
	}

	productPhoto, errPhoto := uc.ProductPhotoService.FindOne(ctx, bson.M{"productId": productID, "key": body.CoverPhoto})
	if errPhoto != nil {
		return entity.InternalServerError(errPhoto.Error())
	}

	coverPhoto := productPhoto.Photo

	// Update data produk
	updateData := bson.M{
		"name":       body.Name,
		"category":   model.AttributeEmbedded{ID: &category.ID, Name: &category.Name},
		"coverPhoto": coverPhoto,
		"price":      body.Price,
		"stock":      body.Stock,
	}

	// Lakukan update produk
	_, err = uc.ProductService.UpdateOne(ctx, bson.M{"_id": productID, "storeId": storeID}, bson.M{"$set": updateData})
	if err != nil {
		return entity.InternalServerError(err.Error())
	}

	return nil
}
