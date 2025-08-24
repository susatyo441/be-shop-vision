package storecontroller

import (
	storedto "be-shop-vision/dto/store"
	usecase "be-shop-vision/usecase/store_usecase"
	"be-shop-vision/util"

	"github.com/gofiber/fiber/v2"
	"github.com/susatyo441/go-ta-utils/response"
)

type makeStoreUseCaseFunc = func() usecase.IStoreUseCase

type StoreController struct {
	UseCase             usecase.IStoreUseCase
	MakeUseCaseFunction makeStoreUseCaseFunc
}

func MakeStoreController(makeUseCaseFunc makeStoreUseCaseFunc) *StoreController {
	return &StoreController{MakeUseCaseFunction: makeUseCaseFunc}
}

// CreateStoregodoc
// @Summary Create Store
// @Description Create Store
// @Tags Store
// @Produce  json
// @Router /api/store [post]
// @Param payload body storedto.CreateStoreDTO true "Payload to create"
// @Security BearerAuth
func (ctrl *StoreController) CreateStore(ctx *fiber.Ctx) error {
	var payload storedto.CreateStoreDTO
	ctx.BodyParser(&payload)
	if err := util.ValidateStruct(payload); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}
	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	err := ctrl.UseCase.CreateStore(ctx.Context(), payload)
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	return response.Created(ctx, "Successfully create store", nil)
}

// GetStores godoc
// @Summary Get Stores
// @Description Get Stores
// @Tags Store
// @Produce  json
// @Router /api/store [get]
// @Security BearerAuth
func (ctrl *StoreController) GetStores(ctx *fiber.Ctx) error {

	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	result, err := ctrl.UseCase.GetStores(ctx.Context())
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	return response.Success(ctx, "Successfully get Stores", result)
}
