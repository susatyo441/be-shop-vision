package productpipeline

import (
	"be-shop-vision/dto"

	"github.com/susatyo441/go-ta-utils/pipeline"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetProductsPipeline(query dto.PaginationQuery, storeID primitive.ObjectID) mongo.Pipeline {
	paginationQuery := pipeline.PaginationQuery{
		Page:      query.Page,
		Limit:     query.Limit,
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
	}

	return pipeline.NewPipelineBuilder().
		Match(bson.M{
			"$or": bson.A{
				bson.M{"name": pipeline.GenerateSearchCondition(query.Search)},
				bson.M{"category.name": pipeline.GenerateSearchCondition(query.Search)},
				bson.M{"stock": pipeline.GenerateSearchCondition(query.Search)},
				bson.M{"price": pipeline.GenerateSearchCondition(query.Search)},
			},
		}).
		Pagination(paginationQuery).
		Build()
}
