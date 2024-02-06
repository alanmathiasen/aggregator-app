package router

import (
	"net/http"

	"github.com/alanmathiasen/aggregator-api/controllers"
	"github.com/alanmathiasen/aggregator-api/middlewares"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func Routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options {
		AllowedOrigins: []string{"http://*", "https://*"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,	
	}))

	
  router.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
		r.Get("/", controllers.GetAllPublicationsHTML)

	})
	
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))
	//--------------------------REST API--------------------------
	// Publications
	router.Get("/api/v1/publications", controllers.GetAllPublications)
	router.Get("/api/v1/publications/{id}", controllers.GetPublicationById)
	router.Post("/api/v1/publications", controllers.CreatePublication)
	router.Put("/api/v1/publications/{id}", controllers.UpdatePublication)
	router.Delete("/api/v1/publications/{id}", controllers.DeletePublication)
	//Chapters
	router.Get("/api/v1/publications/{id}/chapters", controllers.GetAllChaptersByPublicationID)
	router.Post("/api/v1/publications/{id}/chapters", controllers.CreateChapterForPublication)

	router.Get("/auth/login", controllers.LoginHTML)

	router.Get("/{id}", controllers.GetPublicationHTML)
	return router
}