package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github/com/thejhndwn/sanbar/backend/internal/handler/auth"
	"github/com/thejhndwn/sanbar/backend/internal/handler/game"
	"github/com/thejhndwn/sanbar/backend/internal/handler/leaderboard"
	"github/com/thejhndwn/sanbar/backend/internal/handler/user"
)

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", auth.Register).Methods("POST")
	api.HandleFunc("/login", auth.Login).Methods("POST")

	game_api := r.PathPrefix("/api/game").Subrouter()
	game_api.HandleFunc("/{id}/", game.Start)
	game_api.HandleFunc("/{id}/get", game.Get)
	game_api.HandleFunc("/{id}/submit", game.Submit)
	game_api.HandleFunc("/{id}/ready", game.Ready)
	game_api.HandleFunc("/{id}/continue",game.Continue )

	leaderboard_api := r.PathPrefix("/api/leaderboard").Subrouter()
	leaderboard_api.HandleFunc("/survival", leaderboard.Survival).Methods("GET")

	user_api := r.PathPrefix("/user").Subrouter()
	user_api.HandleFunc("/{id}/", user.Profile)
}
