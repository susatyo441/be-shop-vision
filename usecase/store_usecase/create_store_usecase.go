package storeusecase

import (
	storedto "be-shop-vision/dto/store"
	"context"

	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/model"
)

func (uc *StoreUseCase) CreateStore(ctx context.Context, dto storedto.CreateStoreDTO) *entity.HttpError {
	_, err := uc.StoreService.Create(ctx, model.Store{
		Name: dto.Name,
	})
	if err != nil {
		return entity.InternalServerError(err.Error())
	}

	return nil
}
