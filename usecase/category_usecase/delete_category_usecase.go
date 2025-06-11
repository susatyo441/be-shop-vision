package categoryusecase

import (
	"be-shop-vision/dto"
	"context"

	"github.com/susatyo441/go-ta-utils/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc *CategoryUseCase) BulkDeleteCategories(ctx context.Context, body dto.ArrayOfIdDTO, storeID primitive.ObjectID) *entity.HttpError {
	// Konversi slice string menjadi slice primitive.ObjectID
	var objectIDs []primitive.ObjectID
	for _, id := range body.IDs {
		oid, _ := primitive.ObjectIDFromHex(id)

		objectIDs = append(objectIDs, oid)
	}

	// Lakukan penghapusan berdasarkan ObjectID
	deletedCounts, err := uc.CategoryService.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objectIDs}, "storeId": storeID})
	if err != nil {
		return entity.InternalServerError(err.Error())
	}

	if deletedCounts == 0 {
		return entity.BadRequest("Tidak ada kategori yang dihapus, pastikan ID yang diberikan benar")
	}

	return nil
}
