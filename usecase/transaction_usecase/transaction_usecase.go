package transactionusecase

import (
	"be-shop-vision/dto"
	transactiondto "be-shop-vision/dto/transaction"
	"context"

	"github.com/susatyo441/go-ta-utils/db"
	utilDto "github.com/susatyo441/go-ta-utils/dto"
	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/model"
	utilservice "github.com/susatyo441/go-ta-utils/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITransactionUseCase interface {
	CreateTransaction(ctx context.Context, body transactiondto.CreateTransactionDTO, storeID primitive.ObjectID) *entity.HttpError
	GetTransactionList(ctx context.Context, query dto.PaginationQuery, storeID primitive.ObjectID) (*utilDto.PaginationResult[transactiondto.TransactionAggregateDto], *entity.HttpError)
}

type TransactionUseCase struct {
	ProductService     utilservice.Service[model.Product]
	TransactionService utilservice.Service[model.Transaction]
}

func MakeTransactionUseCase() ITransactionUseCase {
	return &TransactionUseCase{
		TransactionService: utilservice.ShopVisionService[model.Transaction](db.TransactionsModelName),
		ProductService:     utilservice.ShopVisionService[model.Product](db.ProductModelName),
	}
}
