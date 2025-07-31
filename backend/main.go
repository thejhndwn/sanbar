package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"context"
	"net/http"
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


	register := RegisterHandler(dbm)
	survival := MakeSurvival(dbm)
	start := Start(dbm)
	get := Get(dbm)
	submit := Submit(dbm)

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", register)
	api.HandleFunc("/login", LoginHandler)
	api.HandleFunc("/newgame", survival)

	game_api := r.PathPrefix("/api/game").Subrouter()
	game_api.HandleFunc("/{id}/", start)
	game_api.HandleFunc("/{id}/get", get)
	game_api.HandleFunc("/{id}/submit", submit)

	leaderboard_api := r.PathPrefix("/api/leaderboard").Subrouter()
	leaderboard_api.HandleFunc("/survival", Survival)

	user_api := r.PathPrefix("/user").Subrouter()
	user_api.HandleFunc("/{id}/", Profile)

	http.ListenAndServe(":8080", r)
}
