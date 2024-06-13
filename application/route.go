package application

import (
	"forum/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func LoadRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Route("/Letrésboncoin", loadOrderRoutes)
	return router

}

func loadOrderRoutes(router chi.Router) {
	orderHandler := &handler.Order{}

	router.Post("/", orderHandler.Create)
	router.Get("/", orderHandler.List)
	router.Get("/{ID}", orderHandler.GetByID)
	router.Put("/{ID}", orderHandler.UpdateByID)
	router.Delete("/{ID}", orderHandler.DeleteByID)
}
