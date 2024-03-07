package main

import (
	"code-runner/util"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/yassinebenaid/godump"
)

func main() {
	// you could insert your favorite logger here for structured or leveled logging

	// load config file app.env
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
	godump.Dump(config)
	switch config.Environment {
	case "development":
		// https://github.com/rs/zerolog?tab=readme-ov-file#pretty-logging
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	case "production":
		gin.SetMode(gin.ReleaseMode)
	}

	// log.Println("Starting server...")
	log.Info().Msg("Starting server...")

	router := gin.Default()

	router.GET("/api/account", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hello": "world",
		})
	})

	srv := &http.Server{
		Addr:    config.ServiceAddress,
		Handler: router,
	}

	// Graceful server shutdown - https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// log.Fatalf("Failed to initialize server: %v\n", err)
			log.Fatal().Err(err).Msg("Failed to initialize server")
		}
	}()

	// log.Printf("Listening on port %v\n", srv.Addr)
	log.Info().Msgf("Listening on port %v\n", srv.Addr)

	// Wait for kill signal of channel
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is passed into the quit channel
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	// log.Println("Shutting down server...")
	log.Info().Msg("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		// log.Fatalf("Server forced to shutdown: %v\n", err)
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}
}
