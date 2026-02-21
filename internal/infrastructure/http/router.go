package http

import (
	"comparei-servico-listas/internal/infrastructure/http/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(handler *ListaHandler) *mux.Router {
	r := mux.NewRouter()

	// Middleware para JSON
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	// --- APLICA O MIDDLEWARE DE AUTH ---
	// Isso protege todas as rotas abaixo
	r.Use(middleware.APIKeyMiddleware)

	// Rotas (agora protegidas)
	r.HandleFunc("/listas", handler.GetListas).Methods("GET")
	r.HandleFunc("/listas", handler.CreateLista).Methods("POST")
	r.HandleFunc("/listas/{id}", handler.GetListaByID).Methods("GET")
	r.HandleFunc("/listas/{id}/finalizar", handler.FinalizarID).Methods("PUT")
	r.HandleFunc("/listas/{id}/itens", handler.AddItem).Methods("POST")
	r.HandleFunc("/listas/{id}/itens", handler.DelItem).Methods("DELETE")
	r.HandleFunc("/itens/{item_id}/check", handler.CheckItem).Methods("PUT")

	return r
}
