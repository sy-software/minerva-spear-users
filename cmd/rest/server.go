package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	minervaLog "github.com/sy-software/minerva-go-utils/log"
	"github.com/sy-software/minerva-spear-users/internal/core/service"
	"github.com/sy-software/minerva-spear-users/internal/handlers"
	"github.com/sy-software/minerva-spear-users/internal/repositories"
)

func main() {
	minervaLog.ConfigureLogger(minervaLog.LogLevel(os.Getenv("LOG_LEVEL")), os.Getenv("CONSOLE_OUTPUT") != "")
	log.Info().Msg("Starting server")

	configRepo := repositories.ConfigRepo{}
	config := configRepo.Get()

	repo := repositories.NewUserRepo(&config)

	authService := service.NewAuthService(repo, config)

	handler := handlers.NewAuthRESTHandler(&config, authService)

	router := gin.Default()

	handler.CreateRoutes(router)

	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Info().Msgf("listen: %s", address)
		} else {
			log.Panic().Err(err).Msg("Can't start server")
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Stack().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
