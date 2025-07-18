package questionerdto

type CreateQuestionerDTO struct {
	Name         string `json:"name" bson:"name"`
	Questioner1  int    `json:"questioner1" bson:"questioner1" validate:"required,gte=1,lte=5"`
	Questioner2  int    `json:"questioner2" bson:"questioner2" validate:"required,gte=1,lte=5"`
	Questioner3  int    `json:"questioner3" bson:"questioner3" validate:"required,gte=1,lte=5"`
	Questioner4  int    `json:"questioner4" bson:"questioner4" validate:"required,gte=1,lte=5"`
	Questioner5  int    `json:"questioner5" bson:"questioner5" validate:"required,gte=1,lte=5"`
	Questioner6  int    `json:"questioner6" bson:"questioner6" validate:"required,gte=1,lte=5"`
	Questioner7  int    `json:"questioner7" bson:"questioner7" validate:"required,gte=1,lte=5"`
	Questioner8  int    `json:"questioner8" bson:"questioner8" validate:"required,gte=1,lte=5"`
	Questioner9  int    `json:"questioner9" bson:"questioner9" validate:"required,gte=1,lte=5"`
	Questioner10 int    `json:"questioner10" bson:"questioner10" validate:"required,gte=1,lte=5"`

	Instagram *string `json:"instagram" bson:"instagram"`
}
