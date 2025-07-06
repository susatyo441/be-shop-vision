package transactiondto

import (
	"time"

	"github.com/susatyo441/go-ta-utils/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionProductAttribute struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	Name string `json:"name" bson:"name"`

	Price int `json:"price" bson:"price"`

	Category model.AttributeEmbedded `json:"category" bson:"category"`

	Quantity int `json:"quantity" bson:"quantity"`

	TotalPrice int `json:"totalPrice" bson:"totalPrice"`

	CoverPhoto string `json:"coverPhoto" bson:"coverPhoto"`
}

type TransactionAggregateDto struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	Products []TransactionProductAttribute `json:"products" bson:"products"`

	TotalPrice int `json:"totalPrice" bson:"totalPrice"`

	StoreID primitive.ObjectID `json:"storeId" bson:"storeId"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type TransactionSummaryDTO struct {
	Monthly              []MonthlySalesDTO `json:"monthly"`
	Daily                []DailySalesDTO   `json:"daily"`
	TotalOrder           int               `json:"totalOrder"`
	TotalOrderThisMonth  int               `json:"totalOrderThisMonth"`
	TotalOrderToday      int               `json:"totalOrderToday"`
	TotalIncome          int               `json:"totalIncome"`
	TotalIncomeThisMonth int               `json:"totalIncomeThisMonth"`
	TotalIncomeToday     int               `json:"totalIncomeToday"`
}

type MonthlySalesDTO struct {
	Month string `json:"month"`
	Sales int    `json:"sales"`
}

type DailySalesDTO struct {
	Day   string `json:"day"`
	Sales int    `json:"sales"`
}
