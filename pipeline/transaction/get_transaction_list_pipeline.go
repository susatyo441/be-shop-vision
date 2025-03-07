package transactionpipeline

import (
	"be-shop-vision/dto"

	"github.com/susatyo441/go-ta-utils/pipeline"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTransactionPipeline(query dto.PaginationQuery, storeID primitive.ObjectID) mongo.Pipeline {
	paginationQuery := pipeline.PaginationQuery{
		Page:      query.Page,
		Limit:     query.Limit,
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
	}

	return pipeline.NewPipelineBuilder().
		Unwind(bson.M{"path": "$products", "preserveNullAndEmptyArrays": true}).

		// Lookup ke products collection untuk ambil coverPhoto
		Lookup(bson.M{
			"from":         "products",
			"localField":   "products._id",
			"foreignField": "_id",
			"as":           "productInfo",
		}).
		// Unwind karena hasil lookup array
		Unwind(bson.M{"path": "$productInfo", "preserveNullAndEmptyArrays": true}).

		// Project semua field yang dibutuhkan
		Project(bson.M{
			"_id":        1,
			"totalPrice": 1,
			"storeId":    1,
			"createdAt":  1,
			"updatedAt":  1,
			"product": bson.M{
				"_id":   "$products._id",
				"name":  "$products.name",
				"price": "$products.price",
				"category": bson.M{
					"_id":  "$products.category._id",
					"name": "$products.category.name",
				},
				"totalPrice": "$products.totalPrice",
				"quantity":   "$products.quantity",
				"coverPhoto": "$productInfo.coverPhoto", // Ambil coverPhoto hasil lookup
			},
		}).

		// Filter berdasarkan nama product atau nama category
		Match(bson.M{
			"$or": bson.A{
				bson.M{"product.name": pipeline.GenerateSearchCondition(query.Search)},
				bson.M{"product.category.name": pipeline.GenerateSearchCondition(query.Search)},
			},
		}).

		// Pagination & Sort
		Pagination(paginationQuery).
		Build()
}
