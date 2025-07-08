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
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	sevenDaysAgo := startOfToday.AddDate(0, 0, -6)

	builder := pipeline.NewPipelineBuilder()

	// Filter berdasarkan store
	builder = builder.Match(bson.M{"storeId": storeID})

	builder = builder.Addfields(bson.M{
		"month": bson.M{"$month": "$createdAt"},
		"day":   bson.M{"$dayOfWeek": "$createdAt"},
	})

	// Sorting data sebelum masuk ke grup agar urutan berdasarkan `createdAt`
	builder = builder.Sort(bson.M{"createdAt": 1})

	// Gunakan facet untuk membuat summary transaksi
	builder = builder.Facet(bson.M{
		// Total transaksi per bulan
		// Total transaksi per bulan
		"monthly": bson.A{
			bson.M{"$match": bson.M{"createdAt": bson.M{"$gte": startOfYear, "$lte": now}}},
			bson.M{"$group": bson.M{
				"_id":   "$month",
				"sales": bson.M{"$sum": "$totalPrice"},
			}},
			bson.M{"$sort": bson.M{"_id": 1}},
			bson.M{"$project": bson.M{
				"month": bson.M{"$switch": bson.M{
					"branches": bson.A{
						bson.M{"case": bson.M{"$eq": bson.A{"$_id", 1}}, "then": "Januari"},
						bson.M{"case": bson.M{"$eq": bson.A{"$_id", 2}}, "then": "Februari"},
						bson.M{"case": bson.M{"$eq": bson.A{"$_id", 3}}, "then": "Maret"},
						bson.M{"case": bson.M{"$eq": bson.A{"$_id", 4}}, "then": "April"},
						bson.M{"case": bson.M{"$eq": bson.A{"$_id", 5}}, "then": "Mei"},
						bson.M{"case": bson.M{"$eq": bson.A{"$_id", 6}}, "then": "Juni"},
						bson.M{"case": bson.M{"$eq": bson.A{"$_id", 7}}, "then": "Juli"},
						bson.M{"case": bson.M{"$eq": bson.A{"$_id", 8}}, "then": "Agustus"},
						bson.M{"case": bson.M{"$eq": bson.A{"$_id", 9}}, "then": "September"},
						bson.M{"case": bson.M{"$eq": bson.A{"$_id", 10}}, "then": "Oktober"},
						bson.M{"case": bson.M{"$eq": bson.A{"$_id", 11}}, "then": "November"},
						bson.M{"case": bson.M{"$eq": bson.A{"$_id", 12}}, "then": "Desember"},
					},
					"default": "Unknown",
				}},
				"sales": 1,
			}},
		},

		// [FIXED] Data harian 7 hari terakhir (termasuk hari ini)
		"daily": bson.A{
			bson.M{"$match": bson.M{
				"createdAt": bson.M{"$gte": sevenDaysAgo, "$lte": now},
			}},
			// Kelompokkan per TANGGAL (bukan hari dalam minggu)
			bson.M{"$group": bson.M{
				"_id": bson.M{
					"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$createdAt"}},
				},
				"sales":       bson.M{"$sum": "$totalPrice"},
				"originalDay": bson.M{"$first": "$day"}, // Simpan hari asli (1-7)
			}},
			// Urutkan dari tanggal TERBARU (descending)
			bson.M{"$sort": bson.M{"_id.date": 1}},
			// Proyeksi: konversi angka hari ke nama hari
			bson.M{"$project": bson.M{
				"day": bson.M{"$switch": bson.M{
					"branches": bson.A{
						bson.M{"case": bson.M{"$eq": bson.A{"$originalDay", 1}}, "then": "Minggu"},
						bson.M{"case": bson.M{"$eq": bson.A{"$originalDay", 2}}, "then": "Senin"},
						bson.M{"case": bson.M{"$eq": bson.A{"$originalDay", 3}}, "then": "Selasa"},
						bson.M{"case": bson.M{"$eq": bson.A{"$originalDay", 4}}, "then": "Rabu"},
						bson.M{"case": bson.M{"$eq": bson.A{"$originalDay", 5}}, "then": "Kamis"},
						bson.M{"case": bson.M{"$eq": bson.A{"$originalDay", 6}}, "then": "Jumat"},
						bson.M{"case": bson.M{"$eq": bson.A{"$originalDay", 7}}, "then": "Sabtu"},
					},
					"default": "Unknown",
				}},
				"sales": 1,
			}},
		},

		// Total seluruh transaksi tanpa batas waktu
		"totalAllTime": bson.A{
			bson.M{"$group": bson.M{
				"_id":         nil,
				"totalOrder":  bson.M{"$sum": 1},
				"totalIncome": bson.M{"$sum": "$totalPrice"},
			}},
		},

		// Total transaksi dalam bulan ini
		"totalThisMonth": bson.A{
			bson.M{"$match": bson.M{"createdAt": bson.M{"$gte": startOfMonth}}},
			bson.M{"$group": bson.M{
				"_id":         nil,
				"totalOrder":  bson.M{"$sum": 1},
				"totalIncome": bson.M{"$sum": "$totalPrice"},
			}},
		},

		// Total transaksi hari ini
		"totalToday": bson.A{
			bson.M{"$match": bson.M{"createdAt": bson.M{"$gte": startOfToday}}},
			bson.M{"$group": bson.M{
				"_id":         nil,
				"totalOrder":  bson.M{"$sum": 1},
				"totalIncome": bson.M{"$sum": "$totalPrice"},
			}},
		},
	})

	// Projeksi akhir agar output lebih clean
	builder = builder.Project(bson.M{
		"monthly": "$monthly",
		"daily":   "$daily",

		"totalOrder":  bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$totalAllTime.totalOrder", 0}}, 0}},
		"totalIncome": bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$totalAllTime.totalIncome", 0}}, 0}},

		"totalOrderThisMonth":  bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$totalThisMonth.totalOrder", 0}}, 0}},
		"totalIncomeThisMonth": bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$totalThisMonth.totalIncome", 0}}, 0}},

		"totalOrderToday":  bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$totalToday.totalOrder", 0}}, 0}},
		"totalIncomeToday": bson.M{"$ifNull": bson.A{bson.M{"$arrayElemAt": bson.A{"$totalToday.totalIncome", 0}}, 0}}},
	)

	return builder.Build()
}
