package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"iss-model-backend/internal/database"
	"iss-model-backend/internal/handlers"
	"iss-model-backend/internal/models"
	"iss-model-backend/internal/services"
)

type Server struct {
	port int

	db          database.Service
	issService  *services.ISSService
	issHandler  *handlers.ISSHandler
	crewService *services.CrewService
	crewHandler *handlers.CrewHandler
	postService *services.PostService
	postHandler *handlers.PostHandler
	authService *services.AuthService
	authHandler *handlers.AuthHandler
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080
	}

	dbService := database.New()

	gormDB := dbService.GetDB()

	err := gormDB.AutoMigrate(&models.ISSPosition{}, &models.Post{}, &models.User{})
	if err != nil {
		fmt.Printf("Failed to auto-migrate models: %v\n", err)
	}

	issService := services.NewISSService(gormDB)
	issHandler := handlers.NewISSHandler(issService)
	crewService := services.NewCrewService()
	crewHandler := handlers.NewCrewHandler(crewService)
	postService := services.NewPostService(gormDB)
	postHandler := handlers.NewPostHandler(postService)
	authService := services.NewAuthService(gormDB)
	authHandler := handlers.NewAuthHandler(authService)

	newServer := &Server{
		port:        port,
		db:          dbService,
		issService:  issService,
		issHandler:  issHandler,
		crewService: crewService,
		crewHandler: crewHandler,
		postService: postService,
		postHandler: postHandler,
		authService: authService,
		authHandler: authHandler,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
