package jsonresponse

import "github.com/gofiber/fiber/v2"

func OK(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": message,
	})
}
