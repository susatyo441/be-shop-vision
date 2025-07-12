package questioncontroller

import (
	questionerdto "be-shop-vision/dto/questioner"
	usecase "be-shop-vision/usecase/questioner_usecase"
	"be-shop-vision/util"

	"github.com/gofiber/fiber/v2"
	"github.com/susatyo441/go-ta-utils/functions"
	"github.com/susatyo441/go-ta-utils/response"
)

type makeQuestionerUseCaseFunc = func() usecase.IQuestionerUseCase

type QuestionerController struct {
	UseCase             usecase.IQuestionerUseCase
	MakeUseCaseFunction makeQuestionerUseCaseFunc
}

func MakeQuestionerController(makeUseCaseFunc makeQuestionerUseCaseFunc) *QuestionerController {
	return &QuestionerController{MakeUseCaseFunction: makeUseCaseFunc}
}

// CreateQuestioner godoc
// @Summary Create Questioner
// @Description Create Questioner
// @Tags Questioner
// @Produce  json
// @Router /api/questioner [post]
// @Param payload body questionerdto.CreateQuestionerDTO true "Payload to create"
// @Security BearerAuth
func (ctrl *QuestionerController) CreateQuestioner(ctx *fiber.Ctx) error {

	var payload questionerdto.CreateQuestionerDTO
	ctx.BodyParser(&payload)
	if err := util.ValidateStruct(payload); err != nil {
		return response.BadRequest(ctx, err.Error(), nil)
	}
	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	err := ctrl.UseCase.CreateQuestioner(ctx.Context(), payload)
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	return response.Created(ctx, "Successfully create questioner", nil)
}

// GetCredits godoc
// @Summary Get Credits list
// @Description Get Credits list
// @Tags Questioner
// @Produce  json
// @Router /api/questioner/credits [get]
// @Security BearerAuth
func (c *QuestionerController) GetCreditList(ctx *fiber.Ctx) error {

	// initialize the use case with logic specific to the company
	c.UseCase = c.MakeUseCaseFunction()

	// get device names using the validated query
	result, err := c.UseCase.GetCredits(ctx.UserContext())
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	// return a success response with the fetched device names and a 200 status
	return response.Success(ctx, "successfully get credit list", result)
}

// GetQuestioner godoc
// @Summary Get Questioner
// @Description Get Questioner
// @Tags Questioner
// @Produce  json
// @Router /api/questioner [get]
// @Security BearerAuth
func (c *QuestionerController) GetQuestioner(ctx *fiber.Ctx) error {

	// initialize the use case with logic specific to the company
	c.UseCase = c.MakeUseCaseFunction()

	// get device names using the validated query
	result, err := c.UseCase.GetQuestioner(ctx.UserContext())
	if err != nil {
		return response.SendResponse(ctx, err.Code, nil, err.Message)
	}

	// return a success response with the fetched device names and a 200 status
	return response.Success(ctx, "successfully get questioner summary", result)
}

// GetQuestionerDetailStats godoc
// @Summary Get Questioner Detail Stats
// @Description Get Questioner Detail Stats
// @Tags Questioner
// @Produce  json
// @Router /api/questioner/stats [get]
// @Security BearerAuth
func (ctrl *QuestionerController) GetQuestionerDetailStats(ctx *fiber.Ctx) error {

	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	questionerDetailStats, errUseCase := ctrl.UseCase.GetQuestionerDetailStats(ctx.Context())
	if errUseCase != nil {
		return response.SendResponse(ctx, errUseCase.Code, nil, errUseCase.Message)
	}

	return response.Success(ctx, "Successfully retrieved questioner detail stats", questionerDetailStats)
}

// GetQuestionerDetail godoc
// @Summary Get Questioner Detail
// @Description Get Questioner Detail
// @Tags Questioner
// @Produce  json
// @Router /api/questioner/{questionerId} [get]
// @Param questionerId path string true "questioner ID"
// @Security BearerAuth
func (ctrl *QuestionerController) GetQuestionerDetail(ctx *fiber.Ctx) error {
	questionerId, err := functions.ParamToObjectID(ctx, "questionerId")
	if err != nil {
		return response.BadRequest(ctx, "Invalid questioner ID format", nil)
	}

	ctrl.UseCase = ctrl.MakeUseCaseFunction()

	questionerDetail, errUseCase := ctrl.UseCase.GetQuestionerDetailByID(ctx.Context(), questionerId)
	if errUseCase != nil {
		return response.SendResponse(ctx, errUseCase.Code, nil, errUseCase.Message)
	}

	return response.Success(ctx, "Successfully retrieved questioner details", questionerDetail)
}
