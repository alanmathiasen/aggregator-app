package router

import (
	"net/http"

	"github.com/alanmathiasen/aggregator-api/auth"
	"github.com/alanmathiasen/aggregator-api/controllers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func Routes() http.Handler {
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
	//Auth
	router.Get("/auth/register", controllers.RegisterPage)
	router.Get("/auth/login", controllers.LoginPage)
	//Loggged in
	router.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)

		r.Get("/discover", controllers.GetAllPublicationsHTML)
		r.Get("/dashboard", controllers.DashboardHTML)
		r.Put("/publication/{id}/follow", controllers.UpsertPublicationFollowHTML)
		r.Delete("/publication/{id}/follow", controllers.DeletePublicationFollowHTML)
	})
	// router.Get("/{id}", controllers.GetPublicationHTML)

	//--------------------------REST API--------------------------
	// Publications
	router.Get("/api/v1/publications", controllers.GetAllPublications)
	// router.Get("/api/v1/publications/{id}", controllers.GetPublicationById)
	router.Post("/api/v1/publications", controllers.CreatePublication)
	router.Put("/api/v1/publications/{id}", controllers.UpdatePublication)
	//router.Delete("/api/v1/publications/{id}", controllers.DeletePublication)
	// Chapters
	router.Get("/api/v1/publications/{id}/chapters", controllers.GetAllChaptersByPublicationID)
	router.Post("/api/v1/publications/{id}/chapters", controllers.CreateChapterForPublication)

	// Auth
	router.Post("/auth/login", controllers.Login)
	router.Post("/auth/register", controllers.Register)
	router.Post("/auth/logout", controllers.Logout)

	return router

}
