package categorycontroller

import (
	"be-shop-vision/dto"
	categorydto "be-shop-vision/dto/category"
	usecase "be-shop-vision/usecase/category_usecase"
	"be-shop-vision/util"

	"github.com/gofiber/fiber/v2"
	"github.com/susatyo441/go-ta-utils/functions"
	"github.com/susatyo441/go-ta-utils/middleware"
	"github.com/susatyo441/go-ta-utils/parser"
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
// @Router /api/category [post]
// @Param payload body categorydto.CreateCategoryDTO true "Payload to create"
// @Security BearerAuth
func (ctrl *CategoryController) CreateCategory(ctx *fiber.Ctx) error {
	var payload categorydto.CreateCategoryDTO
	ctx.BodyParser(&payload)
	if err := util.ValidateStruct(payload); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}
	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	storeId := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)

	err := ctrl.UseCase.CreateCategory(ctx.Context(), payload, storeId)
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	return response.Created(ctx, "Successfully create category", nil)
}

// UpdateCategory godoc
// @Summary Update Category
// @Description Update Category
// @Tags Category
// @Produce  json
// @Router /api/category/{categoryId} [put]
// @Param categoryId path string true "Category ID"
// @Param payload body categorydto.CreateCategoryDTO true "Payload to update"
// @Security BearerAuth
func (ctrl *CategoryController) UpdateCategory(ctx *fiber.Ctx) error {
	var payload categorydto.CreateCategoryDTO
	ctx.BodyParser(&payload)
	if err := util.ValidateStruct(payload); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}

	ctrl.UseCase = ctrl.MakeUseCaseFunction()
	categoryId, paramErr := functions.ParamToObjectID(ctx, "categoryId")
	if paramErr != nil {
		return response.BadRequest(ctx, "Invalid category id format", nil)
	}
	storeId := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)

	err := ctrl.UseCase.UpdateCategory(ctx.Context(), payload, categoryId, storeId)
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	return response.Success(ctx, "Successfully update category", nil)
}

// GetCategoryOptions godoc
// @Summary Get Category Option
// @Description Get Category Option
// @Tags Category
// @Produce  json
// @Router /api/category [get]
// @Param q query dto.PaginationQuery false "Query"
// @Security BearerAuth
func (ctrl *CategoryController) GetCategoryOptions(ctx *fiber.Ctx) error {
	// parse the request query
	rawQuery := ctx.Queries()
	query, _ := parser.ParseQuery[dto.PaginationQuery](rawQuery)

	// validate parsed query
	if err := util.ValidateStruct(query); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}

	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	storeId := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)

	result, err := ctrl.UseCase.GetCategoryOptions(ctx.Context(), storeId, *query)
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	return response.Success(ctx, "Successfully get category options", result)
}

// BulkDeleteCategories godoc
// @Summary Bulk Delete Categories
// @Description Bulk Delete Categories
// @Tags Category
// @Produce  json
// @Router /api/category [delete]
// @Param payload body dto.ArrayOfIdDTO true "Payload to bulk delete"
// @Security BearerAuth
func (ctrl *CategoryController) BulkDeleteCategories(ctx *fiber.Ctx) error {
	var payload dto.ArrayOfIdDTO
	ctx.BodyParser(&payload)
	if err := util.ValidateStruct(payload); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}

	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	storeId := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)

	err := ctrl.UseCase.BulkDeleteCategories(ctx.Context(), payload, storeId)
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	return response.Success(ctx, "Successfully delete categories", nil)
}
