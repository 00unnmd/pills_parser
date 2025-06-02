package transport

import (
	"database/sql"
	"github.com/00unnmd/pills_parser/internal/service"
	"github.com/gorilla/mux"
)

func SetupRouter(db *sql.DB) (*mux.Router, error) {
	authHandler := service.NewAuthHandler(db)
	pillsHandler := service.NewPillsHandler(db)

	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/login", authHandler.Login).Methods("POST")

	pillsRouter := apiRouter.PathPrefix("/pills").Subrouter()
	pillsRouter.Use(authMiddleware)
	pillsRouter.HandleFunc("/Zdravcity", pillsHandler.GetZSPills).Methods("GET")
	pillsRouter.HandleFunc("/Aptekaru", pillsHandler.GetARPills).Methods("GET")
	pillsRouter.HandleFunc("/Eapteka", pillsHandler.GetEAPills).Methods("GET")
	pillsRouter.HandleFunc("/export", pillsHandler.ExportPillsXLSX).Methods("GET")

	return r, nil
}
