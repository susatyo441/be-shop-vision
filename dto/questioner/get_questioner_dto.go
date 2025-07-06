package questionerdto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionerAggregateDto struct {
	Question1 float64 `json:"question1" bson:"question1"`
	Question2 float64 `json:"question2" bson:"question2"`
	Question3 float64 `json:"question3" bson:"question3"`
	Question4 float64 `json:"question4" bson:"question4"`
	Question5 float64 `json:"question5" bson:"question5"`
	Question6 float64 `json:"question6" bson:"question6"`

	Average     float64 `json:"average" bson:"average"`
	TotalAnswer int     `json:"totalAnswer" bson:"totalAnswer"`
}

type QuestionerDetailWithAverage struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Question1 int                `json:"question1" bson:"question1"`
	Question2 int                `json:"question2" bson:"question2"`
	Question3 int                `json:"question3" bson:"question3"`
	Question4 int                `json:"question4" bson:"question4"`
	Question5 int                `json:"question5" bson:"question5"`
	Question6 int                `json:"question6" bson:"question6"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`

	Average float64 `json:"average"`
}
