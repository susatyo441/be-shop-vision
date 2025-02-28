package categoryusecase

import (
	"be-shop-vision/dto"
	categorydto "be-shop-vision/dto/category"
	categorypipeline "be-shop-vision/pipeline/category"
	"context"

	"github.com/susatyo441/go-ta-utils/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (u *CategoryUseCase) GetCategoryOptions(ctx context.Context, storeID primitive.ObjectID, query dto.PaginationQuery) (*categorydto.GetCategoryOptionsResponse, *entity.HttpError) {
	var result []categorydto.GetCategoryOptionsResponse = make([]categorydto.GetCategoryOptionsResponse, 0)

	// perform the aggregation query and handle any errors that occur during the process
	aggregateErr := u.CategoryService.Aggregate(
		&result,
		ctx,
		categorypipeline.GetCategoryOptionPipeline(query, storeID),
		options.Aggregate().SetCollation(&options.Collation{Strength: 3, Locale: "en"}),
	)

	if aggregateErr != nil {
		return nil, entity.InternalServerError(aggregateErr.Error())
	}

	return &result[0], nil
}
