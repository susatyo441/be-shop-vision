package productusecase

import (
	"context"

	"github.com/susatyo441/go-ta-utils/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (uc *ProductUseCase) ExportAllData(ctx context.Context) ([]model.Product, []model.Category, []model.ProductPhoto, []model.Transaction, []model.Store, []model.User, error) {
	products, err := uc.ProductService.Find(ctx, bson.M{})
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	categories, err := uc.CategoryService.Find(ctx, bson.M{})
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	productPhotos, err := uc.ProductPhotoService.Find(ctx, bson.M{})
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	transactions, err := uc.TransactionService.Find(ctx, bson.M{})
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	stores, err := uc.StoreService.Find(ctx, bson.M{})
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	users, err := uc.UserService.Find(ctx, bson.M{})
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return products, categories, productPhotos, transactions, stores, users, nil
}
