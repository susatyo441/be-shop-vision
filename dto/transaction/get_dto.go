package transactiondto

import (
    "time"
    
    "github.com/susatyo441/go-ta-utils/model"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionAggregateDto struct {
	ID         primitive.ObjectID          `json:"_id,omitempty" bson:"_id,omitempty"`
	Product    model.TransactionProductAttribute `json:"product" bson:"product"` 
	TotalPrice int                         `json:"totalPrice" bson:"totalPrice"`
	StoreID    primitive.ObjectID          `json:"storeId" bson:"storeId"`
	CreatedAt  time.Time                   `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time                   `json:"updatedAt" bson:"updatedAt"`
}
