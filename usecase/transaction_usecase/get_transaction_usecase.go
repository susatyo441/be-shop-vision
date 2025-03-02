package transactionusecase

import (
	"be-shop-vision/dto"
	transactionpipeline "be-shop-vision/pipeline/transaction"
	"context"

	utilDto "github.com/susatyo441/go-ta-utils/dto"
	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/functions"
	"github.com/susatyo441/go-ta-utils/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (u *TransactionUseCase) GetTransactionList(ctx context.Context, query dto.PaginationQuery, storeID primitive.ObjectID) (*utilDto.PaginationResult[model.Transaction], *entity.HttpError) {

	var result []utilDto.PaginationResult[model.Transaction]

	// perform the aggregation query and handle any errors that occur during the process
	aggregateErr := u.TransactionService.Aggregate(
		&result,
		ctx,
		transactionpipeline.GetTransactionPipeline(query, storeID),
		options.Aggregate().SetCollation(&options.Collation{Strength: 3, Locale: "en"}),
	)

	if aggregateErr != nil {
		return nil, entity.InternalServerError(aggregateErr.Error())
	}

	response := functions.FormatPaginationResult(result)
	// return the first aggregation result
	return &response, nil
}
