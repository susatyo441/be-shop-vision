package route

import (
	controller "be-shop-vision/controller/product_controller"

	usecase "be-shop-vision/usecase/product_usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/susatyo441/go-ta-utils/middleware"
)

func SetupProductController(app *fiber.App) {
	c := controller.MakeProductController(usecase.MakeProductUseCase)
	// Buat group dengan prefix
	v2 := app.Group("/product")

	// Tambahkan route ke dalam group
	v2.Post("/", middleware.SaveMultiImageMiddleware("product", []string{"firstImage", "secondImage", "thirdImage", "fourthImage", "fifthImage"}, middleware.Medium), c.CreateProduct)
	v2.Put("/:productId", middleware.SaveMultiImageMiddleware("product", []string{"firstImage", "secondImage", "thirdImage", "fourthImage", "fifthImage"}, middleware.Medium), c.UpdateProduct)
	v2.Delete("/", c.BulkDeleteProducts)
	v2.Get("/:productId", c.GetProductDetail)
    v2.Get("/", c.GetProductList)

}
