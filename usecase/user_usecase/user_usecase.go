package userusecase

import (
	dto "be-shop-vision/dto/user"
	"context"
	"mime/multipart"

	"github.com/susatyo441/go-ta-utils/db"
	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/model"
	utilservice "github.com/susatyo441/go-ta-utils/service"
)

type IUserUseCase interface {
	RegisterUser(ctx context.Context, body dto.RegisterUserDTO, files map[string]*multipart.FileHeader) *entity.HttpError
	LoginUser(ctx context.Context, body dto.LoginUserDTO) (*dto.LoginResponseDTO, *entity.HttpError)
}

type UserUseCase struct {
	UserService  utilservice.Service[model.User]
	StoreService utilservice.Service[model.Store]
}

func MakeUserUseCase() IUserUseCase {
	return &UserUseCase{
		UserService:  utilservice.ShopVisionService[model.User](db.UserModelName),
		StoreService: utilservice.ShopVisionService[model.Store](db.StoreModelName),
	}
}
