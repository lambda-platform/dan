package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/dan/controllers"
	"os"
)

func DAN(e *fiber.App) {
	dan := e.Group("/dan")
	dan.Get("/login", controllers.DANRedirect)
	redirectRoute := os.Getenv("DAN_REDIRECT_ROUTE")
	e.Get(redirectRoute, controllers.AuthWithDan)

}
