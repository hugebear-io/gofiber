package fabric

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/hugebear-io/gofiber/errors"
)

type HttpResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

func Response(c *fiber.Ctx, statusCode int, message string, result interface{}) error {
	success := statusCode == http.StatusOK || statusCode == http.StatusCreated || statusCode == http.StatusAccepted
	payload := HttpResponse{
		Success: success,
		Message: message,
		Result:  result,
	}

	return c.Status(statusCode).JSON(payload)
}

func ResponseOK(c *fiber.Ctx, result interface{}) error {
	return Response(c, http.StatusOK, "success", result)
}

func ResponseCreated(c *fiber.Ctx, result interface{}) error {
	return Response(c, http.StatusCreated, "created", result)
}

func ResponseBadRequest(c *fiber.Ctx) error {
	return Response(c, http.StatusBadRequest, "Bad Request", nil)
}

func ResponseUnprocessableEntity(c *fiber.Ctx) error {
	return Response(c, http.StatusUnprocessableEntity, "Unprocessable Entity", nil)
}

func ResponseError(c *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case errors.Error:
		return Response(c, e.StatusCode, e.Message, nil)
	default:
		return Response(c, http.StatusInternalServerError, "System Error", nil)
	}
}
