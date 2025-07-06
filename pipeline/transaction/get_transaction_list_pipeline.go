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

		// Filter awal berdasarkan storeID
		Match(bson.M{"storeId": storeID}).

		// Unwind products agar bisa di-enrich
		Unwind(bson.M{"path": "$products", "preserveNullAndEmptyArrays": true}).

		// Lookup ke collection products untuk ambil coverPhoto
		Lookup(bson.M{
			"from":         "products",
			"localField":   "products._id",
			"foreignField": "_id",
			"as":           "productInfo",
		}).

		// Unwind hasil lookup (productInfo)
		Unwind(bson.M{"path": "$productInfo", "preserveNullAndEmptyArrays": true}).

		// Project field product dengan coverPhoto
		Addfields(bson.M{
			"products.coverPhoto": "$productInfo.coverPhoto",
		}).

		// Group kembali products ke dalam array agar sesuai dengan []TransactionProductAttribute
		Group(bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "totalPrice", Value: bson.D{{Key: "$first", Value: "$totalPrice"}}},
			{Key: "storeId", Value: bson.D{{Key: "$first", Value: "$storeId"}}},
			{Key: "createdAt", Value: bson.D{{Key: "$first", Value: "$createdAt"}}},
			{Key: "updatedAt", Value: bson.D{{Key: "$first", Value: "$updatedAt"}}},
			{Key: "products", Value: bson.D{
				{Key: "$push", Value: bson.D{
					{Key: "_id", Value: "$products._id"},
					{Key: "name", Value: "$products.name"},
					{Key: "price", Value: "$products.price"},
					{Key: "quantity", Value: "$products.quantity"},
					{Key: "totalPrice", Value: "$products.totalPrice"},
					{Key: "coverPhoto", Value: "$products.coverPhoto"},
					{Key: "category", Value: bson.D{
						{Key: "_id", Value: "$products.category._id"},
						{Key: "name", Value: "$products.category.name"},
						{Key: "key", Value: "$products.category.key"},
					}},
				}},
			}},
		}).

		// Optional: filter berdasarkan search query (nama produk atau kategori)
		Match(bson.M{
			"$or": bson.A{
				bson.M{"products.name": pipeline.GenerateSearchCondition(query.Search)},
				bson.M{"products.category.name": pipeline.GenerateSearchCondition(query.Search)},
			},
		}).
		// Pagination & Sort
		Pagination(paginationQuery).
		Build()
}
