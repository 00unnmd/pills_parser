package router

import (
	"database/sql"
	"github.com/00unnmd/pills_parser/handlers"
	"github.com/00unnmd/pills_parser/internals/transport/middleware"
	"github.com/gorilla/mux"
)

func SetupRouter(db *sql.DB) (*mux.Router, error) {
	authHandler := handlers.NewAuthHandler(db)
	pillsHandler := handlers.NewPillsHandler(db)

	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/login", authHandler.Login).Methods("POST")

	pillsRouter := apiRouter.PathPrefix("/pills").Subrouter()
	pillsRouter.Use(middleware.AuthMiddleware)
	pillsRouter.HandleFunc("/Zdravcity", pillsHandler.GetZSPills).Methods("GET")
	pillsRouter.HandleFunc("/Aptekaru", pillsHandler.GetARPills).Methods("GET")
	pillsRouter.HandleFunc("/Eapteka", pillsHandler.GetEAPills).Methods("GET")
	pillsRouter.HandleFunc("/export", pillsHandler.ExportPillsXLSX).Methods("GET")

	return r, nil
}
