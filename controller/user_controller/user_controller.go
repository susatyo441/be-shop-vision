package usercontroller

import (
	userdto "be-shop-vision/dto/user"
	usecase "be-shop-vision/usecase/user_usecase"
	"be-shop-vision/util"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/susatyo441/go-ta-utils/response"
)

type makeUserUseCaseFunc = func() usecase.IUserUseCase

type UserController struct {
	UseCase             usecase.IUserUseCase
	MakeUseCaseFunction makeUserUseCaseFunc
}

func MakeUserController(makeUseCaseFunc makeUserUseCaseFunc) *UserController {
	return &UserController{MakeUseCaseFunction: makeUseCaseFunc}
}

// RegisterUser godoc
// @Summary Register User
// @Description Register User
// @Tags User
// @Produce  json
// @Router /api/user/register [post]
// @Param payload body transactiondto.CreateTransactionDTO true "Payload to create"
// @Security BearerAuth
func (ctrl *UserController) RegisterUser(ctx *fiber.Ctx) error {

	var payload userdto.RegisterUserDTO
	ctx.BodyParser(&payload)
	if err := util.ValidateStruct(payload); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}
	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	files := map[string]*multipart.FileHeader{}
	attributes := "profile_picture"

	fileHeader, err := ctx.FormFile(attributes)
	if err == nil { // File ada
		files[attributes] = fileHeader
	}

	errUseCase := ctrl.UseCase.RegisterUser(ctx.Context(), payload, files)
	if errUseCase != nil {
		return response.SendResponse(ctx, errUseCase.Code, nil, errUseCase.Message)
	}

	return response.Created(ctx, "Successfully register user", nil)
}

// LoginUser godoc
// @Summary Login User
// @Description Login User
// @Tags User
// @Produce  json
// @Router /api/user/login [post]
// @Param payload body userdto.LoginUserDTO true "Payload to login"
// @Security BearerAuth
func (ctrl *UserController) LoginUser(c *fiber.Ctx) error {
	var body userdto.LoginUserDTO

	// Parse request body
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validasi input
	if err := util.ValidateStruct(body); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	// Panggil usecase untuk login
	userData, err := ctrl.UseCase.LoginUser(c.Context(), body)
	if err != nil {
		return response.SendResponse(c, err.Code, nil, err.Message)
	}

	// Simpan session user ID
	c.Locals("session", userData.Token)

	return response.Success(c, "Login berhasil", userData)
}
