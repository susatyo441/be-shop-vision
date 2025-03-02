package route

import (
	controller "be-shop-vision/controller/transaction_controller"

	usecase "be-shop-vision/usecase/transaction_usecase"

	"github.com/gofiber/fiber/v2"
)

func SetupTransactionController(app *fiber.App) {
	c := controller.MakeTransactionController(usecase.MakeTransactionUseCase)
	// Buat group dengan prefix
	v2 := app.Group("/transaction")

	// Tambahkan route ke dalam group
	v2.Post("/", c.CreateTransaction)
	v2.Get("/", c.GetTransactionList)

}
