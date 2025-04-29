// Package main is the entry point of the application.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"wayra/internal/adapter/config"
	"wayra/internal/adapter/repository"
	"wayra/internal/digcontainer"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @title API Specification
// @version 1.0
// @termsOfService http://swagger.io/terms/

// @host localhost:8081
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Provide the Bearer token in the format: "Bearer {token}"

// @tag     name="auth" description="Authentication Endpoints"
// @tag     name="user" description="User Management"
// @tag     name="company" description="Company Management"
// @tag     name="routes" description="Routes Management"
// @tag     name="waypoints" description="Waypoints Management"
// @tag     name="sensor-data" description="Sensor Data Management"
// @tag     name="analytics" description="Analytics Management"
func main() {
	container := digcontainer.BuildContainer()

	err := container.Invoke(func(db *gorm.DB) {
		if err := repository.AutoMigrate(db); err != nil {
			log.Fatalf("Migration failed: %s", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to invoke DB migration: %s", err)
	}

	err = container.Invoke(func(router *gin.Engine, cfg *config.Config) {
		log.Println("Starting server")

		srv := &http.Server{
			Addr:    "localhost:" + strconv.Itoa(cfg.Http.Port),
			Handler: router,
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server failed: %s", err)
			}
		}()

		log.Println("Server is running...")

		<-quit
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %s", err)
		}

		log.Println("Server exiting")
	})

	if err != nil {
		log.Fatalf("Failed to start application: %s", err)
	}
}
