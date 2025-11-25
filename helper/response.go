package helper

import "github.com/gofiber/fiber/v2"

func SuccessResponse(c *fiber.Ctx, status int, data interface{}) error {
    return c.Status(status).JSON(fiber.Map{
        "success": true,
        "data":    data,
    })
}

func ErrorResponse(c *fiber.Ctx, status int, message string, details string) error {
    resp := fiber.Map{
        "success": false,
        "message": message,
    }
    if details != "" {
        resp["error"] = details
    }
    return c.Status(status).JSON(resp)
}
