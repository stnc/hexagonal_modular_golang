package app

import (
	"context"

	"log"
	"syscall"
	"os"
	"time"
	"os/signal"
	"github.com/gofiber/fiber/v3"

)


// func parseID(raw string) (int64, error) {
// 	id, err := strconv.ParseInt(raw, 10, 64)
// 	if err != nil || id <= 0 {
// 		return 0, fmt.Errorf("invalid id")
// 	}
// 	return id, nil
// }

func badRequest(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error":   "bad_request",
		"message": message,
	})
}

func notFound(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error":   "not_found",
		"message": message,
	})
}

func trim(s string) string {
	// Simple trim without importing strings to keep the file compact.
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\n' || s[0] == '\t' || s[0] == '\r') {
		s = s[1:]
	}
	for len(s) > 0 {
		last := s[len(s)-1]
		if last != ' ' && last != '\n' && last != '\t' && last != '\r' {
			break
		}
		s = s[:len(s)-1]
	}
	return s
}


func waitForShutdown(app *fiber.App) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
