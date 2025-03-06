package productusecase

import (
	dto "be-shop-vision/dto/product"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"strings"

	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/functions"
	"github.com/susatyo441/go-ta-utils/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (uc *ProductUseCase) UpdateProduct(ctx context.Context, productID primitive.ObjectID, body dto.UpdateProductDTO, storeID primitive.ObjectID, files map[string]*multipart.FileHeader) *entity.HttpError {
	// 1️⃣ Cari produk lama
	product, err := uc.ProductService.FindOne(ctx, bson.M{"_id": productID, "storeId": storeID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.BadRequest("Produk tidak ditemukan")
		}
		return entity.InternalServerError(err.Error())
	}

	// 2️⃣ Validasi kategori
	categoryObjectID, _ := primitive.ObjectIDFromHex(body.CategoryID)
	category, err := uc.CategoryService.FindOne(ctx, bson.M{"_id": categoryObjectID, "storeId": storeID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.BadRequest("Kategori tidak ditemukan")
		}
		return entity.InternalServerError(err.Error())
	}

	// 3️⃣ Jika nama produk berubah, rename folder + update semua path photo di DB
	if product.Name != body.Name {
		oldFolder := fmt.Sprintf("../acts-files/%s/product/%s", storeID.Hex(), product.Name)
		newFolder := fmt.Sprintf("../acts-files/%s/product/%s", storeID.Hex(), body.Name)

		if err := os.Rename(oldFolder, newFolder); err != nil {
			return entity.InternalServerError("Gagal rename folder produk: " + err.Error())
		}

		// Cari semua product photo
		productPhotos, err := uc.ProductPhotoService.Find(ctx, bson.M{"productId": productID})
		if err != nil {
			return entity.InternalServerError("Gagal mengambil daftar foto produk: " + err.Error())
		}

		// Update setiap path photo di product_photos
		for _, photo := range productPhotos {

			newPath := strings.Replace(photo.Photo, fmt.Sprintf("product/%s", product.Name), fmt.Sprintf("product/%s", body.Name), 1)

			_, err := uc.ProductPhotoService.UpdateOne(ctx, bson.M{"_id": photo.ID}, bson.M{"$set": bson.M{"photo": newPath}})
			if err != nil {
				return entity.InternalServerError("Gagal update path foto produk: " + err.Error())
			}
		}

	}

	// 4️⃣ Simpan gambar baru (jika ada)
	folderName := fmt.Sprintf("product/%s", body.Name)
	attributes := []string{"image1", "image2", "image3", "image4", "image5"}

	imagePaths, _ := functions.SaveMultiImages(storeID, folderName, attributes, files, functions.Medium)

	// 5️⃣ Update product_photos dengan foto baru jika ada
	for key, attribute := range attributes {
		if newPhoto, ok := imagePaths[attribute]; ok {
			oldPhoto, err := uc.ProductPhotoService.FindOneAndUpdate(ctx, bson.M{
				"productId": productID,
				"key":       key + 1,
			}, bson.M{
				"$set": bson.M{"photo": newPhoto},
			}, options.FindOneAndUpdate().SetReturnDocument(options.Before))

			if err != nil {
				return entity.InternalServerError("Gagal update product photo: " + err.Error())
			}

			if err := functions.DeleteImage(fmt.Sprintf("/%s", storeID.Hex()), oldPhoto.Photo); err != nil {
				return entity.InternalServerError("Gagal menghapus foto lama: " + err.Error())
			}

		}
	}

	// 6️⃣ Ambil coverPhoto terbaru
	productPhoto, err := uc.ProductPhotoService.FindOne(ctx, bson.M{
		"productId": productID,
		"key":       body.CoverPhoto,
	})
	if err != nil {
		return entity.InternalServerError("Gagal mendapatkan cover photo: " + err.Error())
	}
	coverPhoto := productPhoto.Photo

	// 7️⃣ Update produk ke database
	updateData := bson.M{
		"name":       body.Name,
		"category":   model.AttributeEmbedded{ID: &category.ID, Name: &category.Name},
		"coverPhoto": coverPhoto,
		"price":      body.Price,
		"stock":      body.Stock,
	}

	if _, err := uc.ProductService.UpdateOne(ctx, bson.M{"_id": productID, "storeId": storeID}, bson.M{"$set": updateData}); err != nil {
		return entity.InternalServerError("Gagal update produk: " + err.Error())
	}

	return nil
}
