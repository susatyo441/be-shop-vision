package productpipeline

import (
	"be-shop-vision/dto"

	"github.com/susatyo441/go-ta-utils/pipeline"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetProductsPipeline(query dto.PaginationQuery, storeID primitive.ObjectID) mongo.Pipeline {
	// Tentukan sort default (createdAt DESC)
	sortBy := "createdAt"
	sortOrder := -1

	// Kalau ada query.SortBy, gunakan itu
	if query.SortBy != "" {
		sortBy = query.SortBy
		sortOrder = query.SortOrder
	}

	paginationQuery := pipeline.PaginationQuery{
		Page:      query.Page,
		Limit:     query.Limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	// Awal pipeline
	builder := pipeline.NewPipelineBuilder()

	// Filter berdasarkan search
	builder = builder.Match(bson.M{
		"$or": bson.A{
			bson.M{"name": pipeline.GenerateSearchCondition(query.Search)},
			bson.M{"category.name": pipeline.GenerateSearchCondition(query.Search)},
			bson.M{"stock": pipeline.GenerateSearchCondition(query.Search)},
			bson.M{"price": pipeline.GenerateSearchCondition(query.Search)},
		},
	})

	// Tambahkan filter jika isAvailable = true
	if query.IsAvailable {
		builder = builder.Match(bson.M{
			"$or": bson.A{
				bson.M{"stock": bson.M{"$gt": 0}},
				bson.M{"variants": bson.M{"$elemMatch": bson.M{"stock": bson.M{"$gt": 0}}}},
			},
		})

		builder = builder.Addfields(bson.M{
			"variants": bson.M{
				"$filter": bson.M{
					"input": "$variants",
					"as":    "variant",
					"cond":  bson.M{"$gt": bson.A{"$$variant.stock", 0}},
				},
			},
		})

	}

	// Tambahkan pagination
	builder = builder.Pagination(paginationQuery)

	// Bangun pipeline
	return builder.Build()
}
