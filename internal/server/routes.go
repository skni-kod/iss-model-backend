// internal/server/routes.go
package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)

	r.Route("/iss", func(r chi.Router) {
		r.Get("/current", s.issHandler.GetCurrentPosition)

		r.Get("/historical/{timestamp}", s.issHandler.GetHistoricalPosition)

		r.Post("/historical", s.issHandler.PostHistoricalRequest)

		r.Get("/range", s.issHandler.GetPositionsInRange)

		r.Get("/status", s.issHandler.GetISSStatus)

		r.Get("/crew", s.crewHandler.GetCurrentCrew)

		r.Get("/crewWithPhotos", s.crewHandler.GetCurrentCrewWithPhotos)
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}

// HelloWorldHandler returns a simple hello world message
// @Summary Hello World
// @Description Returns a hello world message
// @Tags general
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router / [get]
func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"
	resp["service"] = "ISS Model Backend"
	resp["endpoints"] = "/iss/current, /iss/historical/{timestamp}, /iss/range, /iss/status, /iss/crew"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error handling JSON marshal. Err: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResp)
}

// healthHandler returns the health status of the application
// @Summary Health Check
// @Description Returns the health status of the database and application
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(s.db.Health())
	if err != nil {
		log.Printf("error handling JSON marshal. Err: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonResp)
}
