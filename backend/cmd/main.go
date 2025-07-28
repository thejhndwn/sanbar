package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/thejhndwn/sanbar/backend/internal/handler/home"
	"github/com/thejhndwn/sanbar/backend/internal/handler/auth"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home.HomeHandler).Methods("GET")
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", auth.Register).Methods("POST")
	api.HandleFunc("/login", auth.Login).Methods("POST")
}
