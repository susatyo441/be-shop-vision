package route

import (
	controller "be-shop-vision/controller/questioner_controller"
	usecase "be-shop-vision/usecase/questioner_usecase"

	"github.com/gofiber/fiber/v2"
)

func SetupQuestionerController(app *fiber.App) {
	c := controller.MakeQuestionerController(usecase.MakeQuestionerUseCase)
	// Buat group dengan prefix
	v2 := app.Group("/questioner")

	// Tambahkan route ke dalam group
	v2.Post("/", c.CreateQuestioner)
	v2.Get("/credits", c.GetCreditList)
	v2.Get("/stats", c.GetQuestionerDetailStats)
	v2.Get("/:questionerId", c.GetQuestionerDetail)
	v2.Get("/", c.GetQuestioner)

}
