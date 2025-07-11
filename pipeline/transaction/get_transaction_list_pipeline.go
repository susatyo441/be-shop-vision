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

		// --- TAHAP BARU YANG LEBIH ANDAL ---
		// Tambahkan field sementara untuk perhitungan tanggal, lalu buat displayDate

		// Group kembali products ke dalam array
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
		Addfields(bson.M{
			// Hitung awal hari ini di zona waktu 'Asia/Jakarta'
			"todayStart": bson.M{
				"$dateTrunc": bson.M{
					"date":     "$$NOW",
					"unit":     "day",
					"timezone": "Asia/Jakarta",
				},
			},
		}).
		Addfields(bson.M{
			// Hitung awal kemarin berdasarkan awal hari ini
			"yesterdayStart": bson.M{
				"$dateSubtract": bson.M{
					"startDate": "$todayStart",
					"unit":      "day",
					"amount":    1,
				},
			},
		}).

		// --- TAHAP BARU DIMULAI DI SINI ---
		// Tambahkan field 'displayDate' berdasarkan 'createdAt'
		Addfields(bson.M{
			// Sekarang lakukan perbandingan menggunakan variabel yang sudah dihitung di MongoDB
			"displayDate": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$gte": bson.A{"$createdAt", "$todayStart"}},
					"then": "Hari ini",
					"else": bson.M{
						"$cond": bson.M{
							"if":   bson.M{"$gte": bson.A{"$createdAt", "$yesterdayStart"}},
							"then": "Kemarin",
							"else": bson.M{
								// Fallback ke logika format tanggal yang kompleks
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
											}, "default": "",
										},
									}, " ",
									bson.M{"$toString": bson.M{"$year": "$createdAt"}},
								},
							},
						},
					},
				},
			},
		}).

		// Hapus field sementara yang tidak lagi diperlukan
		Project(bson.M{
			"todayStart":     0,
			"yesterdayStart": 0,
		}).
		// --- TAHAP BARU BERAKHIR DI SINI ---

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
