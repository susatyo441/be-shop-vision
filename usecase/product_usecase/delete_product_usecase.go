package productusecase

import (
	"be-shop-vision/dto"
	"context"
	"fmt"
	"os"

	"github.com/susatyo441/go-ta-utils/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (uc *ProductUseCase) BulkDeleteProducts(ctx context.Context, body dto.ArrayOfIdDTO, storeID primitive.ObjectID) *entity.HttpError {
	for _, id := range body.IDs {
		productID, _ := primitive.ObjectIDFromHex(id)

		// 1️⃣ Ambil nama produk sebelum dihapus (diperlukan untuk path folder)
		product, err := uc.ProductService.FindOne(ctx, bson.M{"_id": productID, "storeId": storeID})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				continue // Kalau produk tidak ditemukan, skip (anggap sudah terhapus sebelumnya)
			}
			return entity.InternalServerError(err.Error())
		}

		// 3️⃣ Hapus record di `product_photos` collection
		_, err = uc.ProductPhotoService.DeleteMany(ctx, bson.M{"productId": productID})
		if err != nil {
			return entity.InternalServerError(err.Error())
		}

		// 4️⃣ Hapus record di `products` collection
		_, err = uc.ProductService.DeleteMany(ctx, bson.M{"_id": productID})
		if err != nil {
			return entity.InternalServerError(err.Error())
		}

		// 5️⃣ Hapus folder produk setelah semua data dan file terdaftar dihapus
		productFolderPath := fmt.Sprintf("../acts-files/%s/product/%s", storeID.Hex(), product.Name)
		err = os.RemoveAll(productFolderPath)
		if err != nil {
			return entity.InternalServerError(fmt.Sprintf("Gagal hapus folder produk: %s", err.Error()))
		}
	}

	return nil
}
