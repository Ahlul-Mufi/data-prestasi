package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger(c *fiber.Ctx) error {
    start := time.Now()
    err := c.Next()
    stop := time.Since(start)
    fmt.Printf("[%s] %s %s %d %s\n", start.Format(time.RFC3339), c.Method(), c.Path(), c.Response().StatusCode(), stop)
    return err
}
