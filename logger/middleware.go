package logger

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/hugebear-io/gofiber/fabric"
)

func LoggerMiddleware(c *fiber.Ctx) error {
	if LoggerInstance == nil {
		LoggerInstance = NewLoggerMock()
	}

	logRequestHandler(c)
	err := c.Next()
	logResponseHandler(c)
	return err
}

func logRequestHandler(c *fiber.Ctx) {
	mu := sync.Mutex{}
	mu.Lock()
	query := map[string]interface{}{}
	c.QueryParser(&query)

	body := map[string]interface{}{}
	c.BodyParser(&body)

	LoggerInstance.Debug("",
		LoggerInstance.BuildFields(
			"ip", c.IP(),
			"method", c.Method(),
			"path", c.Path(),
			"query", query,
			"request-body", body,
		)...)
	mu.Unlock()
}

func logResponseHandler(c *fiber.Ctx) {
	code := c.Response().StatusCode()
	body := map[string]interface{}{}
	fabric.Recast(c.Response().Body(), &body)
	LoggerInstance.Debug("",
		LoggerInstance.BuildFields(
			"ip", c.IP(),
			"method", c.Method(),
			"path", c.Path(),
			"status-code", code,
			"response-body", body,
		)...)
}
