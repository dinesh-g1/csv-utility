package main

import (
	"fmt"
	"github.com/dinesh-g1/csv-utility/api"
	"github.com/gorilla/mux"
	"net/http"
)

// Run with
//		go run .
// Send request with:
//		curl -F 'file=@/path/matrix.csv' "localhost:8080/echo"

func main() {
	//mux is the Router
	r := mux.NewRouter()
	v1Api := r.PathPrefix("/api/v1").Subrouter()

	//Registering all the apis
	v1Api.HandleFunc("/echo", api.Echo).Methods("POST")
	v1Api.HandleFunc("/invert", api.Invert).Methods("POST")
	v1Api.HandleFunc("/flatten", api.Flatten).Methods("POST")
	v1Api.HandleFunc("/sum", api.Sum).Methods("POST")
	v1Api.HandleFunc("/multiply", api.Multiply).Methods("POST")

	//Proper logging has to be done
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
	fmt.Println("Stopped the Server")
}
