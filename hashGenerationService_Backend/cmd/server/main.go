package main

import (
	"context"
	"hashGenerationService/internal/handler"
	"hashGenerationService/internal/middleware"
	"hashGenerationService/internal/service"
	"hashGenerationService/internal/store"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	s := store.NewInMemoryStore()
	svc := service.NewService(s)
	h := handler.NewHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /hash", h.GenerateHash)
	mux.HandleFunc("GET /hash/{hash}", h.GetHash)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      middleware.CORS(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Println("Server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
