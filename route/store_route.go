package route

import (
	controller "be-shop-vision/controller/store_controller"

	usecase "be-shop-vision/usecase/store_usecase"

	"github.com/gofiber/fiber/v2"
)

func SetupStoreController(app *fiber.App) {
	c := controller.MakeStoreController(usecase.MakeStoreUseCase)
	// Buat group dengan prefix
	v2 := app.Group("/store")

	// Tambahkan route ke dalam group
	v2.Post("/", c.CreateStore)
	v2.Get("/", c.GetStores)

}
