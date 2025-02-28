package route

import (
	controller "be-shop-vision/controller/category_controller"

	usecase "be-shop-vision/usecase/category_usecase"

	"github.com/gofiber/fiber/v2"
)

func SetupCategoryController(app *fiber.App) {
	c := controller.MakeCategoryController(usecase.MakeCategoryUseCase)
	// Buat group dengan prefix
	v2 := app.Group("/category")

	// Tambahkan route ke dalam group
	v2.Post("/", c.CreateCategory)
	v2.Put("/:categoryId", c.UpdateCategory)
	v2.Get("/", c.GetCategoryOptions)
	v2.Delete("/", c.BulkDeleteCategories)

}
