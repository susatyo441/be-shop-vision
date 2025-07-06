package questionerusecase

import (
	"context"

	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (uc *QuestionerUseCase) GetCredits(ctx context.Context) ([]model.Credit, *entity.HttpError) {
	credits, err := uc.CreditService.Find(ctx, bson.M{})
	if err != nil {
		return nil, entity.InternalServerError(err.Error())
	}

	return credits, nil
}
