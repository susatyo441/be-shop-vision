package categorypipeline

import (
	"be-shop-vision/dto"

	"github.com/susatyo441/go-ta-utils/pipeline"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCategoryOptionPipeline(query dto.PaginationQuery, storeID primitive.ObjectID) mongo.Pipeline {

	return pipeline.NewPipelineBuilder().
		Match(bson.M{
			"storeId": pipeline.GenerateExactFilter(true, storeID),
		}).
		Match(bson.M{
			"$or": bson.A{
				bson.M{"name": pipeline.GenerateSearchCondition(query.Search)},
			},
		}).
		Facet(bson.M{
			"categoryOptions": pipeline.GenerateFacetOption("name", true, "$name", "$_id"),
		}).
		Build()
}
