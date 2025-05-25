/*
Package main provides a simple HTTP API to store and retrieve user credentials and configuration,
as well as download an APK file. It uses the chi router for routing, includes CORS support, and
provides graceful shutdown handling.

Main features:
- POST /set-creds: Save user credentials and associated config
- POST /get-creds: Retrieve saved user configuration using credentials
- GET /get-apk: Download an APK file
*/
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.DefaultLogger)

	r.Post("/get-creds", makeHandler(GetCreds))
	r.Post("/set-creds", makeHandler(SaveCreds))
	r.Get("/get-apk", makeHandler(GetAPK))
	r.Get("/get-server-exe", makeHandler(GetServerExe))
	r.Get("/get-frontend-zip", makeHandler(GetFrontendZip))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Starting server on :8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func makeHandler(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}

		switch e := err.(type) {
		case *APIError:
			log.Printf("[API ERROR] Status: %d, Message: %s", e.Status, e.Msg)
			render.Status(r, e.Status)
			render.JSON(w, r, map[string]any{"error": e.Msg})

		case *NoFieldsFoundError:
			log.Printf("[NOT FOUND] %s", e.Msg)
			render.Status(r, e.Status)
			render.JSON(w, r, map[string]any{"error": e.Msg})

		default:
			log.Printf("[UNKNOWN ERROR] %v", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]any{"error": "internal server error"})
		}
	}
}
