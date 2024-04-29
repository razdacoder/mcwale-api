package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/razdacoder/mcwale-api/services/products"
	"github.com/razdacoder/mcwale-api/services/users"
	"github.com/razdacoder/mcwale-api/utils"
	"gorm.io/gorm"
)

type APIServer struct {
	addr string
	db   *gorm.DB
}

func NewAPISever(addr string, db *gorm.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (server *APIServer) Run() error {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	v1Router := chi.NewRouter()
	v1Router.Get("/status", handleHealth)
	// Users Handlers
	userStore := users.NewStore(server.db)
	userHandler := users.NewHandler(userStore)
	userHandler.RegisterRoutes(v1Router)

	// Product Handlers
	productStore := products.NewStore(server.db)
	productHandler := products.NewHandler(productStore)
	productHandler.RegisterRoutes(v1Router)

	router.Mount("/api/v1", v1Router)
	log.Println("Listening on port ", server.addr)
	return http.ListenAndServe(server.addr, router)
}

func handleHealth(writer http.ResponseWriter, request *http.Request) {
	utils.WriteJSON(writer, http.StatusOK, map[string]string{"status": "OK"})
}
