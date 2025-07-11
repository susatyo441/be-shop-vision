package transactionusecase

import (
	"be-shop-vision/dto"
	transactiondto "be-shop-vision/dto/transaction"
	transactionpipeline "be-shop-vision/pipeline/transaction"
	"context"
	"time"

	utilDto "github.com/susatyo441/go-ta-utils/dto"
	"github.com/susatyo441/go-ta-utils/entity"
	"github.com/susatyo441/go-ta-utils/functions"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (u *TransactionUseCase) GetTransactionList(ctx context.Context, query dto.PaginationQuery, storeID primitive.ObjectID) (*utilDto.PaginationResult[transactiondto.TransactionAggregateDto], *entity.HttpError) {

	var result []utilDto.PaginationResult[transactiondto.TransactionAggregateDto]

	// perform the aggregation query and handle any errors that occur during the process
	aggregateErr := u.TransactionService.Aggregate(
		&result,
		ctx,
		transactionpipeline.GetTransactionPipeline(query, storeID, time.Now()),
		options.Aggregate().SetCollation(&options.Collation{Strength: 3, Locale: "en"}),
	)

	if aggregateErr != nil {
		return nil, entity.InternalServerError(aggregateErr.Error())
	}

	response := functions.FormatPaginationResult(result)
	// return the first aggregation result
	return &response, nil
}

func (u *TransactionUseCase) GetTransactionSummary(ctx context.Context, storeID primitive.ObjectID) (*transactiondto.TransactionSummaryDTO, *entity.HttpError) {

	var result []transactiondto.TransactionSummaryDTO

	// perform the aggregation query and handle any errors that occur during the process
	aggregateErr := u.TransactionService.Aggregate(
		&result,
		ctx,
		transactionpipeline.GetTransactionSummaryPipeline(storeID),
		options.Aggregate().SetCollation(&options.Collation{Strength: 3, Locale: "en"}),
	)

	if aggregateErr != nil {
		return nil, entity.InternalServerError(aggregateErr.Error())
	}

	// return the first aggregation result
	return &result[0], nil
}
