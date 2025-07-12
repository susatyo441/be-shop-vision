package questionerusecase

import (
	dto "be-shop-vision/dto/questioner"
	"context"

	"github.com/susatyo441/go-ta-utils/db"
	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/model"
	utilservice "github.com/susatyo441/go-ta-utils/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IQuestionerUseCase interface {
	CreateQuestioner(ctx context.Context, body dto.CreateQuestionerDTO) *entity.HttpError
	GetCredits(ctx context.Context) ([]model.Credit, *entity.HttpError)
	GetQuestioner(ctx context.Context) (*dto.QuestionerAggregateDto, *entity.HttpError)
	GetQuestionerDetailByID(ctx context.Context, id primitive.ObjectID) (*dto.QuestionerDetailWithAverage, *entity.HttpError)
	GetQuestionerDetailStats(ctx context.Context) (*dto.QuestionerStatistics, *entity.HttpError)
}

type QuestionerUseCase struct {
	QuestionerService utilservice.Service[model.Questioner]
	CreditService     utilservice.Service[model.Credit]
}

func MakeQuestionerUseCase() IQuestionerUseCase {
	return &QuestionerUseCase{
		QuestionerService: utilservice.ShopVisionService[model.Questioner](db.QuestionerModelName),
		CreditService:     utilservice.ShopVisionService[model.Credit](db.CreditModelName),
	}
}
