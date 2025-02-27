package categorycontroller

import (
	dto "be-shop-vision/dto/category"
	usecase "be-shop-vision/usecase/category_usecase"
	"be-shop-vision/util"

	"github.com/gofiber/fiber/v2"
	"github.com/susatyo441/go-ta-utils/middleware"
	"github.com/susatyo441/go-ta-utils/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type makeCategoryUseCaseFunc = func() usecase.ICategoryUseCase

type CategoryController struct {
	UseCase             usecase.ICategoryUseCase
	MakeUseCaseFunction makeCategoryUseCaseFunc
}

func MakeCategoryController(makeUseCaseFunc makeCategoryUseCaseFunc) *CategoryController {
	return &CategoryController{MakeUseCaseFunction: makeUseCaseFunc}
}

// CreateCategory godoc
// @Summary Create Category
// @Description Create Category
// @Tags Category
// @Produce  json
// @Router /category [post]
// @Param payload body dto.CreateCategoryDTO true "Payload to create"
// @Security BearerAuth
func (ctrl *CategoryController) CreateCategory(ctx *fiber.Ctx) error {
	var payload dto.CreateCategoryDTO
	ctx.BodyParser(&payload)
	if err := util.ValidateStruct(payload); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}
	ctrl.UseCase = ctrl.MakeUseCaseFunction()
	userId := ctx.Locals(middleware.UserKey).(primitive.ObjectID)
	storeId := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)

	err := ctrl.UseCase.CreateCategory(ctx.Context(), payload, userId, storeId)
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	return response.Created(ctx, "Successfully create category", nil)
}
