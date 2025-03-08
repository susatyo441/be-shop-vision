package userdto

import "time"

type LoginUserDTO struct {
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required"`
}

type LoginResponseDTO struct {
	Token string          `json:"token"`
	User  UserResponseDTO `json:"user"`
}

type UserResponseDTO struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	StoreID              string    `json:"storeId"`
	PhoneNumber          *string   `json:"phoneNumber,omitempty"`
	Email                string    `json:"email"`
	ProfilePictureSmall  *string   `json:"profilePictureSmall,omitempty"`
	ProfilePictureMedium *string   `json:"profilePictureMedium,omitempty"`
	ProfilePictureBig    *string   `json:"profilePictureBig,omitempty"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}
