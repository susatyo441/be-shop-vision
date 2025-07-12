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
		questioner.Question6 +
		questioner.Question7 + questioner.Question8 + questioner.Question9 + questioner.Question10)

	average := total / 10.0

	// Buat response struct
	detail := &questionerdto.QuestionerDetailWithAverage{
		ID:         questioner.ID,
		Name:       questioner.Name,
		Question1:  questioner.Question1,
		Question2:  questioner.Question2,
		Question3:  questioner.Question3,
		Question4:  questioner.Question4,
		Question5:  questioner.Question5,
		Question6:  questioner.Question6,
		Question7:  questioner.Question7,
		Question8:  questioner.Question8,
		Question9:  questioner.Question9,
		Question10: questioner.Question10,
		CreatedAt:  questioner.CreatedAt,
		UpdatedAt:  questioner.UpdatedAt,
		Average:    average,
	}

	return detail, nil
}

func (uc *QuestionerUseCase) GetQuestionerDetailStats(ctx context.Context) (*questionerdto.QuestionerStatistics, *entity.HttpError) {
	// Ambil semua questioner
	questioners, err := uc.QuestionerService.Find(ctx, bson.M{})
	if err != nil {
		return nil, entity.InternalServerError(err.Error())
	}

	if len(questioners) == 0 {
		return nil, entity.NotFound("Data questioner tidak ditemukan")
	}

	stats := make([]questionerdto.Stat, 10) // 10 pertanyaan

	for _, q := range questioners {
		answers := []int{
			q.Question1, q.Question2, q.Question3, q.Question4, q.Question5,
			q.Question6, q.Question7, q.Question8, q.Question9, q.Question10,
		}

		for i, ans := range answers {
			switch ans {
			case 1:
				stats[i].Count1++
			case 2:
				stats[i].Count2++
			case 3:
				stats[i].Count3++
			case 4:
				stats[i].Count4++
			case 5:
				stats[i].Count5++
			}
			stats[i].Average += float64(ans)
		}
	}

	// Hitung rata-rata
	totalResponden := float64(len(questioners))
	for i := range stats {
		stats[i].Average = stats[i].Average / totalResponden
	}

	// Siapkan response DTO
	response := &questionerdto.QuestionerStatistics{
		TotalRespondents: int(totalResponden),
		Questions:        make([]questionerdto.QuestionStats, 10),
	}

	for i, s := range stats {
		response.Questions[i] = questionerdto.QuestionStats{
			QuestionNumber: i + 1,
			Count1:         s.Count1,
			Count2:         s.Count2,
			Count3:         s.Count3,
			Count4:         s.Count4,
			Count5:         s.Count5,
			Average:        s.Average,
		}
	}

	return response, nil
}
