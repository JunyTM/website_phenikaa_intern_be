package router

import (
	"net/http"
	"phenikaa/controller"
	"phenikaa/infrastructure"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

<<<<<<< HEAD
	// _ "pdt-phenikaa-htdn/docs"
=======
	_ "phenikaa/docs"
>>>>>>> 0766cc6b4f00a86b52515e4a672150604eeb9a9e

	httpSwagger "github.com/swaggo/http-swagger"
)

func Router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(6, "application/json"))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Timeout(time.Duration(5 * time.Second)))
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Use this to allow specific origin hosts
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	// Api swagger for developer mode
	r.Get("/api/v1/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(infrastructure.GetHTTPSwagger()),
		httpSwagger.DocExpansion("none"),
	))

	// Declare controller
	accessController := controller.NewAccessController()

	r.Route("/api/v1", func(router chi.Router) {
		// Public routes
		router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})
		router.Post("/login", accessController.Login)
		
		fs := http.StripPrefix("/api/v1/pnk_intern_storage", http.FileServer(http.Dir(infrastructure.GetRootPath()+"/"+infrastructure.GetStoragePath())))
		router.Get("/pnk_intern_storage/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fs.ServeHTTP(w, r)
		}))
	})
	return r
}
