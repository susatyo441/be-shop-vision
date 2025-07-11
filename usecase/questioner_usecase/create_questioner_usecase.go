package questionerusecase

import (
	dto "be-shop-vision/dto/questioner"
	"context"

	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/model"
)

func (uc *QuestionerUseCase) CreateQuestioner(ctx context.Context, body dto.CreateQuestionerDTO) *entity.HttpError {
	questioner, err := uc.QuestionerService.Create(ctx, model.Questioner{
		Name:       body.Name,
		Question1:  body.Questioner1,
		Question2:  body.Questioner2,
		Question3:  body.Questioner3,
		Question4:  body.Questioner4,
		Question5:  body.Questioner5,
		Question6:  body.Questioner6,
		Question7:  body.Questioner7,
		Question8:  body.Questioner8,
		Question9:  body.Questioner9,
		Question10: body.Questioner10,
	})
	if err != nil {
		return entity.InternalServerError(err.Error())
	}
	if body.Instagram != nil && *body.Instagram != "" {
		_, err = uc.CreditService.Create(ctx, model.Credit{
			QuestionerID: questioner.ID,
			Instagram:    *body.Instagram,
			Name:         questioner.Name,
		})
		if err != nil {
			return entity.InternalServerError(err.Error())
		}
	}

	return nil
}
