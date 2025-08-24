package storedto

type CreateStoreDTO struct {
	Name string `json:"name" bson:"name" validate:"required"`
}
