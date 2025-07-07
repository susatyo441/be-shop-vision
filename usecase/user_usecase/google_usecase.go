package userusecase

import (
	dto "be-shop-vision/dto/user"
	"be-shop-vision/util"
	"context"

	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/functions"
	"github.com/susatyo441/go-ta-utils/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	googleOauth2 "google.golang.org/api/oauth2/v2"
)

// 2. Implementasikan metode baru di UserUseCase
func (uc *UserUseCase) LoginGoogleCallback(ctx context.Context, googleUserInfo *googleOauth2.Userinfo) (*dto.LoginResponseDTO, *entity.HttpError) {
	// Cari user berdasarkan email yang didapat dari Google
	user, err := uc.UserService.FindOne(ctx, bson.M{"email": googleUserInfo.Email})

	// Jika user tidak ditemukan, buat user baru
	if err != nil {
		if err == mongo.ErrNoDocuments {
			storeId, _ := primitive.ObjectIDFromHex("67cd1f507488199eeb9b8579")
			// Buat user baru dari data Google
			newUser := model.User{
				ID:    primitive.NewObjectID(),
				Name:  googleUserInfo.Name,
				Email: googleUserInfo.Email,
				Store: model.AttributeEmbedded{ID: functions.MakePointer(storeId), Name: functions.MakePointer("Warung Mamah Tyas")},
			}

			_, errCreate := uc.UserService.Create(ctx, newUser)
			// Buat user baru di DB
			if errCreate != nil {
				return nil, entity.InternalServerError(errCreate.Error())
			}
			// Gunakan user yang baru dibuat
			user = &newUser
		} else {
			// Error lain saat query ke DB
			return nil, entity.InternalServerError(err.Error())
		}
	}

	// ---- BUAT JWT TOKEN ----
	// Pada titik ini, 'user' adalah user yang sudah ada atau yang baru dibuat.
	// Sekarang, kita buat token JWT untuknya.
	token, errToken := util.GenerateJWT(user.ID.Hex(), user.Store.ID.Hex(), false) // Sesuaikan dengan fungsi JWT Anda
	if errToken != nil {
		return nil, entity.InternalServerError("Gagal membuat token")
	}

	// Kembalikan response login yang sukses
	return &dto.LoginResponseDTO{
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
	}, nil
}
