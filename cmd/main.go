package main

import (
	"flag"
	"github.com/dinesh-g1/csv-utility/api"
	"github.com/dinesh-g1/csv-utility/consts"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	host := flag.String("host", "localhost", "host address")
	port := flag.String("port", "8080", "port number")
	flag.Parse()

	r := mux.NewRouter()
	v1Api := r.PathPrefix("/api/v1").Subrouter()

	//Registering all the api endpoints
	v1Api.Handle("/echo", api.RootHandler(api.Echo)).Methods("POST")
	v1Api.Handle("/invert", api.RootHandler(api.Invert)).Methods("POST")
	v1Api.Handle("/flatten", api.RootHandler(api.Flatten)).Methods("POST")
	v1Api.Handle("/sum", api.RootHandler(api.Sum)).Methods("POST")
	v1Api.Handle("/multiply", api.RootHandler(api.Multiply)).Methods("POST")

	log.Printf("Server started on port %s\n", *port)

	addr := *host + consts.COLON + *port
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}

	log.Printf("Server stopped on port %s\n", *port)
}
