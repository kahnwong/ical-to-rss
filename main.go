package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	slogfiber "github.com/samber/slog-fiber"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

func main() {
	mode := os.Getenv("MODE")
	logger, err := newLogger()
	if err != nil {
		logger.Error("Invalid log level", "value", os.Getenv("LOG_LEVEL"), "error", err)
		os.Exit(1)
	}

	slog.SetDefault(logger)

	var listenAddress string
	switch mode {
	case "PRODUCTION":
		listenAddress = ":3000"
	case "DEVELOPMENT":
		listenAddress = "localhost:3000"
	default:
		logger.Error("Listen address is not set")
		os.Exit(1)
	}

	if err := newApp(logger).Listen(listenAddress); err != nil {
		logger.Error("Fiber app error", "error", err)
		os.Exit(1)
	}
}

func newApp(logger *slog.Logger) *fiber.App {
	app := fiber.New()
	app.Use(slogfiber.New(logger))

	// 60 requests per 1 minute max per client IP
	app.Use(limiter.New(limiter.Config{
		Max:        60,
		Expiration: time.Minute,
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).SendString(
				"Too many requests. Try again in " + c.GetRespHeader(fiber.HeaderRetryAfter) + "s",
			)
		},
	}))

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("ICAL to RSS")
	})
	app.Get("/feed", FeedHandler(logger))

	return app
}

func newLogger() (*slog.Logger, error) {
	level := slog.LevelDebug
	var err error
	if value := os.Getenv("LOG_LEVEL"); value != "" {
		err = level.UnmarshalText([]byte(value))
	}

	zerologger := zerolog.New(os.Stderr)
	if os.Getenv("MODE") == "DEVELOPMENT" {
		zerologger = zerologger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	logger := slog.New(slogzerolog.Option{
		Level:  level,
		Logger: &zerologger,
	}.NewZerologHandler())
	return logger, err
}
