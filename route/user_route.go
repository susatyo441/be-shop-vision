package route

import (
	controller "be-shop-vision/controller/user_controller"
	usecase "be-shop-vision/usecase/user_usecase"

	"github.com/gofiber/fiber/v2"
)

func SetupUserController(app *fiber.App) {
	c := controller.MakeUserController(usecase.MakeUserUseCase)
	// Buat group dengan prefix
	v2 := app.Group("/user")

	// Tambahkan route ke dalam group
	v2.Post("/register", c.RegisterUser)
	v2.Post("/login", c.LoginUser)

}
