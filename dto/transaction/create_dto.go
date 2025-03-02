package transactiondto

type CreateTransactionDTO struct {
	ProductID string `json:"productID" bson:"productID" validate:"required"`
	Quantity  int    `json:"quantity" bson:"quantity" validate:"required,gte=1"`
}
