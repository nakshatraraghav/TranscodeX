package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func UserRouter(router *chi.Mux) {

	subrouter := chi.NewRouter()
	router.Mount("/users", subrouter)

	subrouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}
