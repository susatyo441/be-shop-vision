package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/susatyo441/go-ta-utils/middleware"
)

func SetupRoute(app *fiber.App) {
	SetupUserController(app)
	app.Use(middleware.ValidateJWT())

	SetupCategoryController(app)
	SetupStoreController(app)
	SetupProductController(app)
	SetupTransactionController(app)
	SetupQuestionerController(app)
}
