package transactioncontroller

import (
	"be-shop-vision/dto"
	transactiondto "be-shop-vision/dto/transaction"
	usecase "be-shop-vision/usecase/transaction_usecase"
	"be-shop-vision/util"

	"github.com/gofiber/fiber/v2"
	"github.com/susatyo441/go-ta-utils/middleware"
	"github.com/susatyo441/go-ta-utils/parser"
	"github.com/susatyo441/go-ta-utils/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type makeTransactionUseCaseFunc = func() usecase.ITransactionUseCase

type TransactionController struct {
	UseCase             usecase.ITransactionUseCase
	MakeUseCaseFunction makeTransactionUseCaseFunc
}

func MakeTransactionController(makeUseCaseFunc makeTransactionUseCaseFunc) *TransactionController {
	return &TransactionController{MakeUseCaseFunction: makeUseCaseFunc}
}

// CreateTransaction godoc
// @Summary Create Transaction
// @Description Create Transaction
// @Tags Transaction
// @Produce  json
// @Router /api/transaction [post]
// @Param payload body transactiondto.CreateTransactionDTO true "Payload to create"
// @Security BearerAuth
func (ctrl *TransactionController) CreateTransaction(ctx *fiber.Ctx) error {

	var payload transactiondto.CreateTransactionDTO
	ctx.BodyParser(&payload)
	if err := util.ValidateStruct(payload); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}
	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	storeId := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)

	err := ctrl.UseCase.CreateTransaction(ctx.Context(), payload, storeId)
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	return response.Created(ctx, "Successfully create transaction", nil)
}

// GetTransactionList godoc
// @Summary Get Transaction list
// @Description Get Transaction list
// @Tags Transaction
// @Produce  json
// @Router /api/transaction [get]
// @Param q query dto.PaginationQuery false "Query"
// @Security BearerAuth
func (c *TransactionController) GetTransactionList(ctx *fiber.Ctx) error {

	rawQuery := ctx.Queries()
	query, _ := parser.ParseQuery[dto.PaginationQuery](rawQuery)

	// validate parsed query
	if err := util.ValidateStruct(query); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}

	storeID := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)

	// initialize the use case with logic specific to the company
	c.UseCase = c.MakeUseCaseFunction()

	// get device names using the validated query
	result, err := c.UseCase.GetTransactionList(ctx.UserContext(), *query, storeID)
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	// return a success response with the fetched device names and a 200 status
	return response.Success(ctx, "successfully get transaction list", result)
}

// GetTransactionSummary godoc
// @Summary Get Summary list
// @Description Get Summary list
// @Tags Transaction
// @Produce  json
// @Router /api/transaction/summary [get]
// @Security BearerAuth
func (c *TransactionController) GetTransactionSummary(ctx *fiber.Ctx) error {

	storeID := ctx.Locals(middleware.StoreKey).(primitive.ObjectID)

	// initialize the use case with logic specific to the company
	c.UseCase = c.MakeUseCaseFunction()

	// get device names using the validated query
	result, err := c.UseCase.GetTransactionSummary(ctx.UserContext(), storeID)
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	// return a success response with the fetched device names and a 200 status
	return response.Success(ctx, "successfully get transaction summary", result)
}
