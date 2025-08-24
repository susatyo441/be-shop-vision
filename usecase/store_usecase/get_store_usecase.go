package storeusecase

import (
	"context"

	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (uc *StoreUseCase) GetStores(ctx context.Context) ([]model.Store, *entity.HttpError) {
	stores, err := uc.StoreService.Find(ctx, bson.M{})
	if err != nil {
		return nil, entity.InternalServerError(err.Error())
	}

	return stores, nil
}
