package util

import "github.com/gofiber/fiber/v2"

type ApiResponse struct {
	Status     string      `json:"status"`
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message,omitempty"`
	StatusCode int         `json:"status_code"`
}

func HttpError(ctx *fiber.Ctx, statusCode int, message string) error {
	return ctx.Status(statusCode).JSON(ApiResponse{
		Status:     "fail",
		StatusCode: statusCode,
		Message:    message,
	})
}

func HttpSuccess(ctx *fiber.Ctx, statusCode int, data interface{}) error {
	return ctx.Status(statusCode).JSON(ApiResponse{
		Status:     "success",
		StatusCode: statusCode,
		Data:       data,
	})
}
