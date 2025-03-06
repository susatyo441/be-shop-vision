package transactiondto

type CreateTransactionDTO struct {
	Data []CreateTransactionAttrDTO `json:"data" bson:"data" validate:"required,dive,required"`
}

type CreateTransactionAttrDTO struct {
	ProductID string `json:"productID" bson:"productID" validate:"required"`
	Quantity  int    `json:"quantity" bson:"quantity" validate:"required,gte=1"`
}
