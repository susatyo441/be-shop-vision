package transactionpipeline

import (
	"time"

	"github.com/susatyo441/go-ta-utils/pipeline"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTransactionSummaryPipeline(storeID primitive.ObjectID) mongo.Pipeline {
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	builder := pipeline.NewPipelineBuilder()

	// Filter store
	builder = builder.Match(bson.M{"storeId": storeID})

	// Add fields month & dayOfWeek (numeric)
	builder = builder.Addfields(bson.M{
		"month": bson.M{"$dateToString": bson.M{"format": "%b", "date": "$createdAt"}},
		"day": bson.M{
			"$switch": bson.M{
				"branches": bson.A{
					bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$dayOfWeek": "$createdAt"}, 1}}, "then": "Minggu"},
					bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$dayOfWeek": "$createdAt"}, 2}}, "then": "Senin"},
					bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$dayOfWeek": "$createdAt"}, 3}}, "then": "Selasa"},
					bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$dayOfWeek": "$createdAt"}, 4}}, "then": "Rabu"},
					bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$dayOfWeek": "$createdAt"}, 5}}, "then": "Kamis"},
					bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$dayOfWeek": "$createdAt"}, 6}}, "then": "Jumat"},
					bson.M{"case": bson.M{"$eq": bson.A{bson.M{"$dayOfWeek": "$createdAt"}, 7}}, "then": "Sabtu"},
				},
				"default": "Unknown",
			},
		},
	})

	// Facet - pecah hasil ke banyak bagian
	builder = builder.Facet(bson.M{
		"monthly": bson.A{
			bson.M{"$group": bson.M{
				"_id":   "$month",
				"sales": bson.M{"$sum": "$totalPrice"},
			}},
			bson.M{"$project": bson.M{
				"month": "$_id",
				"sales": 1,
			}},
		},

		"daily": bson.A{
			bson.M{"$group": bson.M{
				"_id":   "$day",
				"sales": bson.M{"$sum": "$totalPrice"},
			}},
			bson.M{"$project": bson.M{
				"day":   "$_id",
				"sales": 1,
			}},
		},

		"total": bson.A{
			bson.M{"$group": bson.M{
				"_id":         nil,
				"totalOrder":  bson.M{"$sum": 1},
				"totalIncome": bson.M{"$sum": "$totalPrice"},
			}},
		},

		"totalThisMonth": bson.A{
			bson.M{"$match": bson.M{"createdAt": bson.M{"$gte": startOfMonth}}},
			bson.M{"$group": bson.M{
				"_id":         nil,
				"totalOrder":  bson.M{"$sum": 1},
				"totalIncome": bson.M{"$sum": "$totalPrice"},
			}},
		},

		"totalToday": bson.A{
			bson.M{"$match": bson.M{"createdAt": bson.M{"$gte": startOfToday}}},
			bson.M{"$group": bson.M{
				"_id":         nil,
				"totalOrder":  bson.M{"$sum": 1},
				"totalIncome": bson.M{"$sum": "$totalPrice"},
			}},
		},
	})

	// Final projection (flatten semua hasil facet)
	builder = builder.Project(bson.M{
		"monthly": "$monthly",
		"daily":   "$daily",

		"totalOrder":  bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$total.totalOrder", 0}}, 0}},
		"totalIncome": bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$total.totalIncome", 0}}, 0}},

		"totalOrderThisMonth":  bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$totalThisMonth.totalOrder", 0}}, 0}},
		"totalIncomeThisMonth": bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$totalThisMonth.totalIncome", 0}}, 0}},

		"totalOrderToday":  bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$totalToday.totalOrder", 0}}, 0}},
		"totalIncomeToday": bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$totalToday.totalIncome", 0}}, 0}},
	})

	return builder.Build()
}
