//go:generate swagger generate spec
// Package main RedFishLine.
//
// the purpose of this application is to provide an interface for ILO Configuration through API
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes:
//     Host:
//     BasePath:
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Julien SENON <julien.senon@gmail.com>

// Example server https://16.248.65.69/redfish/v1/

package main

import (
	"api"
	// "fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"web"
)

// TO FIX

func main() {
	r := mux.NewRouter()

	// Remove CORS Header check to allow swagger and application on same host and port
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	// To be changed
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "PATCH"})

	// Web Part
	r.HandleFunc("/index", web.Index)
	r.HandleFunc("/help", web.Help)
	r.HandleFunc("/debug", web.Debug)

	// Static dir
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("templates/static/"))))

	// Health Check
	r.HandleFunc("/healthy/am-i-up", api.Statusamiup).Methods("GET")
	r.HandleFunc("/healthy/about", api.Statusabout).Methods("GET")

	http.ListenAndServe(":9040", handlers.CORS(originsOk, headersOk, methodsOk)(r))
}
