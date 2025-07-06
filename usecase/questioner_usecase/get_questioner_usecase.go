package questionerusecase

import (
	questionerdto "be-shop-vision/dto/questioner"
	questionpipeline "be-shop-vision/pipeline/question"
	"context"
	"errors"

	"github.com/susatyo441/go-ta-utils/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (u *QuestionerUseCase) GetQuestioner(ctx context.Context) (*questionerdto.QuestionerAggregateDto, *entity.HttpError) {

	var result []questionerdto.QuestionerAggregateDto

	// perform the aggregation query and handle any errors that occur during the process
	aggregateErr := u.QuestionerService.Aggregate(
		&result,
		ctx,
		questionpipeline.GetQuestionerAggregatePipeline(),
		options.Aggregate().SetCollation(&options.Collation{Strength: 3, Locale: "en"}),
	)

	if aggregateErr != nil {
		return nil, entity.InternalServerError(aggregateErr.Error())
	}

	return &result[0], nil
}

func (uc *QuestionerUseCase) GetQuestionerDetailByID(ctx context.Context, id primitive.ObjectID) (*questionerdto.QuestionerDetailWithAverage, *entity.HttpError) {
	// Ambil detail by ID
	questioner, err := uc.QuestionerService.FindOne(ctx, bson.M{"_id": id})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entity.NotFound("Questioner tidak ditemukan")
		}
		return nil, entity.InternalServerError(err.Error())
	}

	// Hitung average (float64)
	total := float64(questioner.Question1 +
		questioner.Question2 +
		questioner.Question3 +
		questioner.Question4 +
		questioner.Question5 +
		questioner.Question6)

	average := total / 6.0

	// Buat response struct
	detail := &questionerdto.QuestionerDetailWithAverage{
		ID:        questioner.ID,
		Name:      questioner.Name,
		Question1: questioner.Question1,
		Question2: questioner.Question2,
		Question3: questioner.Question3,
		Question4: questioner.Question4,
		Question5: questioner.Question5,
		Question6: questioner.Question6,
		CreatedAt: questioner.CreatedAt,
		UpdatedAt: questioner.UpdatedAt,
		Average:   average,
	}

	return detail, nil
}
