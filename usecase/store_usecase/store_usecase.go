package storeusecase

import (
	storedto "be-shop-vision/dto/store"
	"context"

	"github.com/susatyo441/go-ta-utils/db"
	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/model"
	utilservice "github.com/susatyo441/go-ta-utils/service"
)

type IStoreUseCase interface {
	CreateStore(ctx context.Context, body storedto.CreateStoreDTO) *entity.HttpError
	GetStores(ctx context.Context) ([]model.Store, *entity.HttpError)
}

type StoreUseCase struct {
	StoreService utilservice.Service[model.Store]
}

func MakeStoreUseCase() IStoreUseCase {
	return &StoreUseCase{

		StoreService: utilservice.ShopVisionService[model.Store](db.StoreModelName),
	}
}
