package main

import (
	"hashGenerationService/internal/handler"
	"hashGenerationService/internal/middleware"
	"hashGenerationService/internal/service"
	"hashGenerationService/internal/store"
	"log"
	"net/http"
)

func main() {

	//	 DI

	s := store.NewInMemoryStore()
	svc := service.NewService(s)
	h := handler.NewHandler(svc)

	// mux is request router
	mux := http.NewServeMux()
	// POST
	mux.HandleFunc("/hash", h.GenerateHash)

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", middleware.CORS(mux)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
