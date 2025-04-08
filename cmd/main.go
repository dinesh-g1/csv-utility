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
	host := flag.String(consts.KeyHost, consts.DefaultHost, "host address")
	port := flag.String(consts.KeyPort, consts.DefaultPort, "port number")
	flag.Parse()

	r := mux.NewRouter()

	//Versioning the api
	v1Api := r.PathPrefix(consts.ApiV1).Subrouter()

	//Registering all the api endpoints
	v1Api.Handle(consts.EndpointEcho, api.RootHandler(api.Echo)).Methods(http.MethodPost)
	v1Api.Handle(consts.EndpointInvert, api.RootHandler(api.Invert)).Methods(http.MethodPost)
	v1Api.Handle(consts.EndpointFlatten, api.RootHandler(api.Flatten)).Methods(http.MethodPost)
	v1Api.Handle(consts.EndpointSum, api.RootHandler(api.Sum)).Methods(http.MethodPost)
	v1Api.Handle(consts.EndpointMultiply, api.RootHandler(api.Multiply)).Methods(http.MethodPost)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf("Server started on port %s\n", *port)

	addr := *host + consts.Colon + *port
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}

	log.Printf("Server stopped on port %s\n", *port)
}
