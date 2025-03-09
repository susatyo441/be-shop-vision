package productusecase

import (
	dto "be-shop-vision/dto/product"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/functions"
	"github.com/susatyo441/go-ta-utils/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (uc *ProductUseCase) CreateProduct(ctx context.Context, body dto.CreateProductDTO, storeID primitive.ObjectID, files map[string]*multipart.FileHeader) *entity.HttpError {
	categoryObjectID, _ := primitive.ObjectIDFromHex(body.CategoryID)

	// Validasi coverPhoto dalam rentang index yang benar
	if body.CoverPhoto > 5 {
		return entity.BadRequest("Index coverPhoto tidak valid")
	}

	category, err := uc.CategoryService.FindOne(ctx, bson.M{"_id": categoryObjectID, "storeId": storeID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.BadRequest("Kategori tidak ditemukan")
		}
		return entity.InternalServerError(err.Error())
	}

	_, err = uc.ProductService.FindOne(ctx, bson.M{
		"name": body.Name, "storeId": storeID,
	})

	if err == nil {
		return entity.BadRequest("Nama produk duplikat")
	} else if err != mongo.ErrNoDocuments {
		return entity.InternalServerError(err.Error())
	}

	attributes := []string{"image1", "image2", "image3", "image4", "image5"}

	// Simpan semua image ke filesystem (ini di usecase)
	folderName := fmt.Sprintf("product/%s", body.Name)
	imagePaths, err := functions.SaveMultiImages(storeID, folderName, attributes, files, functions.Medium)
	if err != nil {
		return entity.InternalServerError("Gagal menyimpan gambar: " + err.Error())
	}

	// Ubah map ke slice biar urut
	var photos []string
	for _, attr := range attributes {
		if path, ok := imagePaths[attr]; ok {
			photos = append(photos, path)
		} else {
			return entity.InternalServerError(fmt.Sprintf("Gambar %s gagal disimpan", attr))
		}
	}

	coverPhoto := photos[body.CoverPhoto-1]
	variants := make([]model.ProductVariantsAttr, len(body.Variants))
	for i, variantDTO := range body.Variants {
		variants[i] = model.ProductVariantsAttr{
			Name:         variantDTO.Name,
			Price:        variantDTO.Price,
			CapitalPrice: nil,
			Stock:        variantDTO.Stock,
		}
	}
	product := model.Product{
		ID:         primitive.NewObjectID(),
		Name:       body.Name,
		StoreID:    storeID,
		Category:   model.AttributeEmbedded{ID: &category.ID, Name: &category.Name},
		Stock:      body.Stock,
		CoverPhoto: coverPhoto,
		Price:      body.Price,
		Variants:   variants,
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
