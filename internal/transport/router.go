package transport

import (
	"database/sql"
	"github.com/00unnmd/pills_parser/internal/service"
	"github.com/gorilla/mux"
)

func SetupRouter(db *sql.DB) (*mux.Router, error) {
	authHandler := service.NewAuthHandler(db)
	pillsHandler := service.NewPillsHandler(db)
	optionsHandler := service.NewOptionsHandler(db)

	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/login", authHandler.Login).Methods("POST")

	pillsRouter := apiRouter.PathPrefix("/pills").Subrouter()
	pillsRouter.Use(authMiddleware)
	pillsRouter.HandleFunc("/ozon", pillsHandler.GetOzonPills).Methods("GET")
	pillsRouter.HandleFunc("/mnn", pillsHandler.GetMNNPills).Methods("GET")
	pillsRouter.HandleFunc("/competitors", pillsHandler.GetCompetitorsPills).Methods("GET")
	pillsRouter.HandleFunc("/export", pillsHandler.ExportPillsXLSX).Methods("GET")
	pillsRouter.HandleFunc("/options", optionsHandler.GetOptions).Methods("GET")

	return r, nil
}
