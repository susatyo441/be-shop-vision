package userusecase

import (
	dto "be-shop-vision/dto/user"
	"context"
	"mime/multipart"

	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/functions"
	"github.com/susatyo441/go-ta-utils/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (uc *UserUseCase) RegisterUser(ctx context.Context, body dto.RegisterUserDTO, files map[string]*multipart.FileHeader) *entity.HttpError {
	// Cek apakah email sudah terdaftar
	existingUser, _ := uc.UserService.FindOne(ctx, bson.M{"email": body.Email})
	if existingUser != nil {
		return entity.BadRequest("Email sudah digunakan")
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.InternalServerError("Gagal mengenkripsi password")
	}

	// Konversi StoreID menjadi ObjectID
	storeID, err := primitive.ObjectIDFromHex(body.StoreID)
	if err != nil {
		return entity.BadRequest("StoreID tidak valid")
	}

	store, errStore := uc.StoreService.FindOne(ctx, bson.M{"_id": storeID})

	if errStore != nil {
		return entity.InternalServerError(errStore.Error())
	}

	folderName := "profile_picture"
	attributes := []string{"profile_picture"}

	lowImage, errPhoto := functions.SaveMultiImages(storeID, folderName, attributes, files, functions.Low)

	mediumImage, errPhoto := functions.SaveMultiImages(storeID, folderName, attributes, files, functions.Medium)

	highImage, errPhoto := functions.SaveMultiImages(storeID, folderName, attributes, files, functions.High)

	if errPhoto != nil {
		return entity.InternalServerError("Gagal menyimpan gambar: " + errPhoto.Error())
	}

	// Buat user baru
	newUser := model.User{
		ID:                   primitive.NewObjectID(),
		Name:                 body.Name,
		Store:                model.AttributeEmbedded{ID: &storeID, Name: &store.Name},
		PhoneNumber:          body.PhoneNumber,
		Email:                body.Email,
		ProfilePictureSmall:  functions.MakePointer(lowImage[folderName]),
		ProfilePictureMedium: functions.MakePointer(mediumImage[folderName]),
		ProfilePictureBig:    functions.MakePointer(highImage[folderName]),
		Password:             functions.MakePointer(string(hashedPassword)),
	}

	// Simpan ke database
	_, err = uc.UserService.Create(ctx, newUser)
	if err != nil {
		return entity.InternalServerError("Gagal menyimpan user")
	}

	return nil
}
