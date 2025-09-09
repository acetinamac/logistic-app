// Package main Logistics API.
// @title           Logistics API
// @version         1.0
// @description     API del backend de logística.
// @contact.name    Equipo Backend
// @contact.email   devs@example.com
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"net/http"
	"os"

	"logistics-app/backend/internal/app"

	_ "logistics-app/backend/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	r := mux.NewRouter()
	if err := app.Bootstrap(r); err != nil {
		log.Fatal(err)
	}

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	addr := ":8080"
	if v := os.Getenv("API_PORT"); v != "" {
		addr = ":" + v
	}
	log.Printf("API listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
