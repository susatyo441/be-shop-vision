package questionerpipeline

import (
	"github.com/susatyo441/go-ta-utils/pipeline"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetQuestionerAggregatePipeline() mongo.Pipeline {
	return pipeline.NewPipelineBuilder().
		// Group: hitung rata-rata per pertanyaan & total answer
		Group(bson.D{
			{Key: "_id", Value: nil},
			{Key: "question1", Value: bson.D{{Key: "$avg", Value: "$question1"}}},
			{Key: "question2", Value: bson.D{{Key: "$avg", Value: "$question2"}}},
			{Key: "question3", Value: bson.D{{Key: "$avg", Value: "$question3"}}},
			{Key: "question4", Value: bson.D{{Key: "$avg", Value: "$question4"}}},
			{Key: "question5", Value: bson.D{{Key: "$avg", Value: "$question5"}}},
			{Key: "question6", Value: bson.D{{Key: "$avg", Value: "$question6"}}},
			{Key: "totalAnswer", Value: bson.D{{Key: "$sum", Value: 1}}},
		}).
		// Tambahkan average keseluruhan
		Addfields(bson.M{
			"average": bson.M{
				"$avg": bson.A{
					"$question1",
					"$question2",
					"$question3",
					"$question4",
					"$question5",
					"$question6",
				},
			},
		}).
		// Project sesuai DTO
		Project(bson.M{
			"_id":         0,
			"question1":   1,
			"question2":   1,
			"question3":   1,
			"question4":   1,
			"question5":   1,
			"question6":   1,
			"average":     1,
			"totalAnswer": 1,
		}).
		Build()
}
