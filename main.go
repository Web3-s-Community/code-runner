package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"code-runner/util"

	ginalgorand "code-runner/module/algorand/transport/gin"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/yassinebenaid/godump"
)

func main() {

	// load config file app.env
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
	switch config.Environment {
	case "development":
		godump.Dump(config)
	case "production":
		gin.SetMode(gin.ReleaseMode)

	}

	log.Logger = log.Output(util.LoggerOutput(config))

	// log.Println("Starting server...")
	log.Info().Msg("Starting server...")

	router := gin.Default()

	// router.Use(util.CORSMiddleware())
	router.Use(util.CORSConfig())

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello Server Compiler!!")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	codeV1 := router.Group("/code/v1")
	{
		codeV1.POST("/algo-playground", ginalgorand.ExecuteCodePlaygroundHandler)
	}

	if config.CertFilePath != "" && config.KeyFilePath != "" {
		router.RunTLS(config.ServiceAddress, config.CertFilePath, config.KeyFilePath)
	}

	srv := &http.Server{
		Addr:    config.ServiceAddress,
		Handler: router,
	}

	// Graceful server shutdown - https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/close/server.go
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// log.Fatalf("Failed to initialize server: %v\n", err)
			log.Fatal().Err(err).Msg("Failed to initialize server")
		}
	}()

	// log.Printf("Listening on port %v\n", srv.Addr)
	log.Info().Msgf("Listening on port %v", srv.Addr)

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
