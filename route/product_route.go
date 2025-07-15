package route

import (
	controller "be-shop-vision/controller/product_controller"
	usecase "be-shop-vision/usecase/product_usecase"

	"github.com/gofiber/fiber/v2"
)

func SetupProductController(app *fiber.App) {
	c := controller.MakeProductController(usecase.MakeProductUseCase)
	app.Get("/export-all", c.ExportAll)
	// Buat group dengan prefix
	v2 := app.Group("/product")

	// Tambahkan route ke dalam group
	v2.Post("/", c.CreateProduct)
	v2.Put("/:productId", c.UpdateProduct)
	v2.Get("/export-photos", c.ExportAllProductPhotos)
	v2.Delete("/", c.BulkDeleteProducts)
	v2.Get("/:productId", c.GetProductDetail)
	v2.Get("/", c.GetProductList)

}
