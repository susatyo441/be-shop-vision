package productdto

type CreateProductDTO struct {
	CategoryID string `json:"categoryID" bson:"categoryID" validate:"required"`
	Name       string `json:"name" bson:"name" validate:"required"`
	Price      int    `json:"price" bson:"price" validate:"required,gte=0"`
	CoverPhoto int    `json:"coverPhoto" bson:"coverPhoto" validate:"required,gte=1,lte=5"`
	Stock      int    `json:"stock" bson:"stock" validate:"required,gte=0"`
}

type ProductPhotosDTO struct {
	Image1 string `json:"image1,omitempty"`
	Image2 string `json:"image2,omitempty"`
	Image3 string `json:"image3,omitempty"`
	Image4 string `json:"image4,omitempty"`
	Image5 string `json:"image5,omitempty"`
}

type UpdateProductDTO struct {
	CategoryID string           `json:"categoryID" bson:"categoryID" validate:"required"`
	Name       string           `json:"name" bson:"name" validate:"required"`
	CoverPhoto int              `json:"coverPhoto" bson:"coverPhoto" validate:"required,gte=1,lte=5"`
	Photos     ProductPhotosDTO `json:"photos"`
	Price      int              `json:"price" bson:"price" validate:"required,gte=0"`
	Stock      int              `json:"stock" bson:"stock" validate:"required,gte=0"`
}
