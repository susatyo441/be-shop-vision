package productdto

type ProductVariantDTO struct {
	Name  string `json:"name" bson:"name" validate:"required"`
	Price int    `json:"price" bson:"price" validate:"required,gte=1"`
	Stock int    `json:"stock" bson:"stock" validate:"required,gte=0"`
}

type CreateProductDTO struct {
	CategoryID string              `json:"categoryID" bson:"categoryID" validate:"required"`
	Name       string              `json:"name" bson:"name" validate:"required"`
	Price      *int                `json:"price" bson:"price" validate:"omitempty,gte=1"`
	CoverPhoto int                 `json:"coverPhoto" bson:"coverPhoto" validate:"required,gte=1,lte=5"`
	Stock      *int                `json:"stock" bson:"stock" validate:"omitempty,gte=0"`
	Variants   []ProductVariantDTO `json:"variants" bson:"variants"`
}

type UpdateProductDTO struct {
	CategoryID string              `json:"categoryID" bson:"categoryID" validate:"required"`
	Name       string              `json:"name" bson:"name" validate:"required"`
	CoverPhoto int                 `json:"coverPhoto" bson:"coverPhoto" validate:"required,gte=1,lte=5"`
	Price      *int                `json:"price" bson:"price" validate:"omitempty,gte=1"`
	Stock      *int                `json:"stock" bson:"stock" validate:"omitempty,gte=0"`
	Variants   []ProductVariantDTO `json:"variants" bson:"variants"`
}

type UpdateProductStockDTO struct {
	Stock    int     `json:"stock" bson:"stock" validate:"required,gte=0"`
	Variants *string `json:"variants" bson:"variants" validate:"omitempty"`
}
