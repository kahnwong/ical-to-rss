package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	// init
	mode := os.Getenv("MODE")

	// ------------ api server ------------ //
	// entrypoint
	listenAddress := ""
	switch mode {
	case "PRODUCTION":
		gin.SetMode(gin.ReleaseMode)
		listenAddress = ":3000"
	case "DEVELOPMENT":
		gin.SetMode(gin.DebugMode)
		listenAddress = "localhost:3000"
	default:
		log.Fatal().Msg("Listen address is not set")
	}

	// app
	app := gin.New()

	// set logger
	isPrettyLog := false
	if mode == "DEVELOPMENT" {
		isPrettyLog = true
	}
	zerologger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	if isPrettyLog {
		zerologger = zerologger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	app.Use(logger.SetLogger(logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
		return zerologger
	})))

	// 60 requests per 1 minute max
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  60,
	}
	store := memory.NewStore()
	rateLimiter := limiter.New(store, rate)
	app.Use(mgin.NewMiddleware(rateLimiter))

	// routes
	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ICAL to RSS")
	})
	app.GET("/feed", FeedHandler(zerologger))

	// start server
	if err := app.Run(listenAddress); err != nil {
		zerologger.Fatal().Err(err).Msg("Gin app error")
	}
}
