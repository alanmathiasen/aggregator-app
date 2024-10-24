package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alanmathiasen/aggregator-api/internal/auth"
	"github.com/alanmathiasen/aggregator-api/internal/db"
	"github.com/alanmathiasen/aggregator-api/internal/handlers"
	"github.com/alanmathiasen/aggregator-api/internal/services"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

type Application struct {
	Config Config
	Models services.Models
}

func (app *Application) Serve() error {
	port := os.Getenv("PORT")
	fmt.Println("Server running on port", port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.RegisterRoutes(),
	}
	return srv.ListenAndServe()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	auth.InitStore()

	cfg := Config{
		Port: os.Getenv(("PORT")),
	}

	dsn := os.Getenv("DSN")
	dbConn, err := db.ConnectPostgres(dsn)
	if err != nil {
		fmt.Println("Couldn't open DB")
	}

	defer dbConn.DB.Close()

	if err != nil {
		fmt.Println("Couldn't connect to DB")
	}

	app := &Application{
		Config: cfg,
		Models: services.New(dbConn.DB),
	}

	log.Fatal(app.Serve())
}

func (app *Application) RegisterRoutes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Use(auth.SessionMiddleware)

	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// -------------------------HTML SERVER------------------------
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)
	})
	// Auth
	router.Get("/auth/register", handlers.RegisterPage)
	router.Get("/auth/login", handlers.LoginPage)
	// Loggged in
	router.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		r.Get("/discover", handlers.GetAllPublicationsHTML)
		r.Get("/dashboard", handlers.DashboardHTML)
		r.Get("/publication/{id}", handlers.GetPublicationHTML)
		// r.Put("/publication/{id}/follow", handlers.UpsertPublicationFollowHTML)
		// r.Delete("/publication/{id}/follow", handlers.DeletePublicationFollowHTML)
	})
	// router.Get("/{id}", handlers.GetPublicationHTML)

	//--------------------------REST API--------------------------
	// Publications
	// router.Get("/api/v1/publications", handlers.GetAllPublications)
	// router.Get("/api/v1/publications/{id}", handlers.GetPublicationById)
	router.Post("/api/v1/publications", handlers.CreatePublication)
	router.Put("/api/v1/publications/{id}", handlers.UpdatePublication)
	// router.Delete("/api/v1/publications/{id}", handlers.DeletePublication)
	// Chapters
	router.Get("/api/v1/publications/{id}/chapters", handlers.GetAllChaptersByPublicationID)
	router.Post("/api/v1/publications/{id}/chapters", handlers.CreateChapterForPublication)

	// Auth
	router.Post("/auth/login", handlers.Login)
	router.Post("/auth/register", handlers.Register)
	router.Post("/auth/logout", handlers.Logout)

	return router
}
