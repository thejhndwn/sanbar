package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"context"
)

func main()	{
	fmt.Println("entered main")
	r := mux.NewRouter()

	dbc := GetConfigFromEnv()
	dbm := NewDatabaseManager(dbc)
	c := context.Background()

	if err := dbm.Initialize(c); err != nil {
		fmt.Printf("DATABSE INITAILIXE NOT GOOD ABORT: %s", err)
		return 
	}




	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", RegisterHandler)
	api.HandleFunc("/login", LoginHandler)

	game_api := r.PathPrefix("/api/game").Subrouter()
	game_api.HandleFunc("/{id}/", Start)
	game_api.HandleFunc("/{id}/get", Get)
	game_api.HandleFunc("/{id}/submit", Submit)
	game_api.HandleFunc("/{id}/ready", Ready)

	leaderboard_api := r.PathPrefix("/api/leaderboard").Subrouter()
	leaderboard_api.HandleFunc("/survival", Survival)

	user_api := r.PathPrefix("/user").Subrouter()
	user_api.HandleFunc("/{id}/", Profile)
}
