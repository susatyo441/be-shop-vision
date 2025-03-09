package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	_ "be-shop-vision/docs"
	"be-shop-vision/route"
)

// InitializeApp initializes the Fiber app with routes and middleware.
// Used for test
func InitializeApp(isTest bool) (*fiber.App, string) {
	err := godotenv.Load(".env")
	if err != nil && !isTest {
		log.Fatal("Error loading .env file:", err)
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, // 100MB
	})

	// BasicAuth middleware configuration
	authConfig := basicauth.Config{
		Users: map[string]string{
			os.Getenv("DOCS_USERNAME"): os.Getenv("DOCS_PASSWORD"),
		},
		Realm: "Restricted",
	}

	// Apply BasicAuth to '/docs*' route
	app.Use("/docs/*", basicauth.New(authConfig))
	app.Use(cors.New())

	// Set Swagger route to /assets/docs
	app.Get("/docs/*", swagger.New(swagger.Config{
		URL: "/docs/doc.json", // Specify the path to your OpenAPI file, delete /asset-control if in local
	}))

	// Setup route config
	route.SetupRoute(app)

	port := os.Getenv("PORT")

	//Return for test
	return app, port
}

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @BasePath /
func main() {
	app, port := InitializeApp(false)
	app.Listen(":" + port)
}
