package userdto

type RegisterUserDTO struct {
	Name        string  `json:"name" bson:"name" validate:"required"`
	StoreID     string  `json:"storeID" bson:"storeID" validate:"required"`
	PhoneNumber *string `json:"phoneNumber" bson:"phoneNumber" validate:"omitempty"`
	Email       string  `json:"email" bson:"email" validate:"required,email"`
	Password    string  `json:"password" bson:"password" validate:"required,min=6"`
}
