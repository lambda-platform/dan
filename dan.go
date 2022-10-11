package dan

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/dan/routes"
)

func Set(e *fiber.App) {
	routes.DAN(e)
}
