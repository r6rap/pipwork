package main

import (
	"log"
	"net/http"

	"github.com/r6rap/pipwork/internal/db"
	"github.com/r6rap/pipwork/internal/api"
	"github.com/gorilla/mux"
)

func main() {
	db.ConnectMySQL()

	r := mux.NewRouter()

	r.HandleFunc("/api/logs", api.GetLogs).Methods("GET")

	r.HandleFunc("/api/targets", api.GetTarget).Methods("GET")
	r.HandleFunc("/api/target/{id}", api.GetTargetByID).Methods("GET")
	r.HandleFunc("/api/targets", api.CreateTarget).Methods("POST")
	r.HandleFunc("/api/targets/{id}", api.UpdateTarget).Methods("PUT")
	r.HandleFunc("/api/targets/{id}", api.DeleteTarget).Methods("DELETE")

	log.Println("🌐 Server running in :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}