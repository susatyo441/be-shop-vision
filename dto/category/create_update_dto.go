package categorydto

type CreateCategoryDTO struct {
	Name string `json:"name" bson:"name" validate:"required"`
}
