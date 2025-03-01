package productusecase

import (
	"be-shop-vision/dto"
	"context"
	"fmt"

	"github.com/susatyo441/go-ta-utils/entity"
    "github.com/susatyo441/go-ta-utils/functions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc *ProductUseCase) BulkDeleteProducts(ctx context.Context, body dto.ArrayOfIdDTO, storeID primitive.ObjectID) *entity.HttpError {
	for _, id := range body.IDs {
		productID, _ := primitive.ObjectIDFromHex(id)
		productPhotos, err := uc.ProductPhotoService.Find(ctx, bson.M{"productId": productID})
		if err != nil {
			return entity.InternalServerError(err.Error())
		}

		for _, productPhoto := range productPhotos {
			err = functions.DeleteImage(fmt.Sprintf("/%s", storeID.Hex()), productPhoto.Photo)
			if err != nil {
				return entity.InternalServerError(err.Error())
			}
		}

		_, err = uc.ProductPhotoService.DeleteMany(ctx, bson.M{"productId": productID})
		if err != nil {
			return entity.InternalServerError(err.Error())
		}

		_, err = uc.ProductService.DeleteMany(ctx, bson.M{"_id": productID})
		if err != nil {
			return entity.InternalServerError(err.Error())
		}

	}

	return nil
}
