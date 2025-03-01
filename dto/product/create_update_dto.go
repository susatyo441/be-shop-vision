package productdto

type CreateProductDTO struct {
	CategoryID string `json:"categoryID" bson:"categoryID" validate:"required"`
	Name       string `json:"name" bson:"name" validate:"required"`
	Price      int    `json:"price" bson:"price" validate:"required, gte=0"`
	CoverPhoto int    `json:"coverPhoto" bson:"coverPhoto" validate:"required,gte=1,lte=5"`
}

type ProductPhotosDTO struct {
	FirstImage  string `json:"firstImage,omitempty"`
	SecondImage string `json:"secondImage,omitempty"`
	ThirdImage  string `json:"thirdImage,omitempty"`
	FourthImage string `json:"fourthImage,omitempty"`
	FifthImage  string `json:"fifthImage,omitempty"`
}

type UpdateProductDTO struct {
	CategoryID string           `json:"categoryID" bson:"categoryID" validate:"required"`
	Name       string           `json:"name" bson:"name" validate:"required"`
	CoverPhoto int              `json:"coverPhoto" bson:"coverPhoto" validate:"required,gte=1,lte=5"`
	Photos     ProductPhotosDTO `json:"photos"`
	Price      int              `json:"price" bson:"price" validate:"required, gte=0"`
}
