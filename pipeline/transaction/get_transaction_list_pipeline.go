package transactionpipeline

import (
	"be-shop-vision/dto"
	"time"

	"github.com/susatyo441/go-ta-utils/pipeline"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTransactionPipeline(query dto.PaginationQuery, storeID primitive.ObjectID, now time.Time) mongo.Pipeline {
	paginationQuery := pipeline.PaginationQuery{
		Page:      query.Page,
		Limit:     query.Limit,
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
	}

	// Hitung waktu di zona waktu Asia/Jakarta
	loc, _ := time.LoadLocation("Asia/Jakarta")
	nowInJakarta := now.In(loc)
	todayStart := time.Date(nowInJakarta.Year(), nowInJakarta.Month(), nowInJakarta.Day(), 0, 0, 0, 0, loc)
	yesterdayStart := todayStart.AddDate(0, 0, -1)

	// Konversi ke tipe MongoDB
	todayStartPrim := primitive.NewDateTimeFromTime(todayStart)
	yesterdayStartPrim := primitive.NewDateTimeFromTime(yesterdayStart)

	// Ekspresi untuk displayDate
	displayDateExpr := bson.M{
		"$cond": bson.M{
			"if":   bson.M{"$gte": bson.A{"$createdAt", todayStartPrim}},
			"then": "Hari ini",
			"else": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$gte": bson.A{"$createdAt", yesterdayStartPrim}},
					"then": "Kemarin",
					"else": bson.M{
						"$concat": bson.A{
							bson.M{"$toString": bson.M{"$dayOfMonth": "$createdAt"}}, " ",
							bson.M{
								"$switch": bson.M{
									"branches": bson.A{
										bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$month": "$createdAt"}, 1}}, "then": "Januari"},
										bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$month": "$createdAt"}, 2}}, "then": "Februari"},
										bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$month": "$createdAt"}, 3}}, "then": "Maret"},
										bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$month": "$createdAt"}, 4}}, "then": "April"},
										bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$month": "$createdAt"}, 5}}, "then": "Mei"},
										bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$month": "$createdAt"}, 6}}, "then": "Juni"},
										bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$month": "$createdAt"}, 7}}, "then": "Juli"},
										bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$month": "$createdAt"}, 8}}, "then": "Agustus"},
										bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$month": "$createdAt"}, 9}}, "then": "September"},
										bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$month": "$createdAt"}, 10}}, "then": "Oktober"},
										bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$month": "$createdAt"}, 11}}, "then": "November"},
										bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$month": "$createdAt"}, 12}}, "then": "Desember"},
									},
									"default": "",
								},
							},
							" ",
							bson.M{"$toString": bson.M{"$year": "$createdAt"}},
						},
					},
				},
			},
		},
	}

	return pipeline.NewPipelineBuilder().
		Match(bson.M{"storeId": storeID}).
		Unwind(bson.M{"path": "$products", "preserveNullAndEmptyArrays": true}).
		Lookup(bson.M{
			"from":         "products",
			"localField":   "products._id",
			"foreignField": "_id",
			"as":           "productInfo",
		}).
		Unwind(bson.M{"path": "$productInfo", "preserveNullAndEmptyArrays": true}).
		Addfields(bson.M{
			"products.coverPhoto": "$productInfo.coverPhoto",
		}).
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
		Addfields(bson.M{"displayDate": displayDateExpr}). // Gunakan ekspresi yang sudah dihitung
		Match(bson.M{
			"$or": bson.A{
				bson.M{"products.name": pipeline.GenerateSearchCondition(query.Search)},
				bson.M{"products.category.name": pipeline.GenerateSearchCondition(query.Search)},
			},
		}).
		Pagination(paginationQuery).
		Build()
}