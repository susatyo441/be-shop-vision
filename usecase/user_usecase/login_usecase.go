package userusecase

import (
	dto "be-shop-vision/dto/user"
	"be-shop-vision/util"
	"context"

	"github.com/susatyo441/go-ta-utils/entity"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func (uc *UserUseCase) LoginUser(ctx context.Context, body dto.LoginUserDTO) (*dto.LoginResponseDTO, *entity.HttpError) {
	// Cari user berdasarkan email
	user, err := uc.UserService.FindOne(ctx, bson.M{"email": body.Email})
	if err != nil {
		return nil, entity.BadRequest("Email atau password salah")
	}

	// Cek apakah password cocok
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(body.Password)); err != nil {
		return nil, entity.BadRequest("Email atau password salah")
	}

	// Generate JWT token dengan storeId & userId
	token, err := util.GenerateJWT(user.ID.Hex(), user.Store.ID.Hex(), false)
	if err != nil {
		return nil, entity.InternalServerError("Gagal membuat token")
	}

	// Buat response tanpa password
	response := &dto.LoginResponseDTO{
		Token: token,
		User: dto.UserResponseDTO{
			ID:                   user.ID.Hex(),
			Name:                 user.Name,
			StoreID:              user.Store.ID.Hex(),
			PhoneNumber:          user.PhoneNumber,
			Email:                user.Email,
			ProfilePictureSmall:  user.ProfilePictureSmall,
			ProfilePictureMedium: user.ProfilePictureMedium,
			ProfilePictureBig:    user.ProfilePictureBig,
			CreatedAt:            user.CreatedAt,
			UpdatedAt:            user.UpdatedAt,
		},
	}

	return response, nil
}
