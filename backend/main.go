package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"context"
	"net/http"
)

func main()	{
	fmt.Println("entered main")
	r := mux.NewRouter()

	dbc := GetConfigFromEnv()
	dbm := NewDatabaseManager(dbc)
	c := context.Background()

   cors_middle := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"}, // Frontend origin
        AllowedMethods: []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type"},
    })

	if err := dbm.Initialize(c); err != nil {
		fmt.Printf("DATABSE INITAILIXE NOT GOOD ABORT: %s", err)
		return 
	}

	login := LoginHandler(dbm)
	register := RegisterHandler(dbm)
	survival := MakeSurvival(dbm)
	start := Start(dbm)
	get := Get(dbm)
	submit := Submit(dbm)

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", register)
	api.HandleFunc("/login", login)
	api.HandleFunc("/newgame", survival)

	game_api := r.PathPrefix("/api/game").Subrouter()
	game_api.HandleFunc("/{id}/", start)
	game_api.HandleFunc("/{id}/get", get)
	game_api.HandleFunc("/{id}/submit", submit)

	leaderboard_api := r.PathPrefix("/api/leaderboard").Subrouter()
	leaderboard_api.HandleFunc("/survival", Survival)

	user_api := r.PathPrefix("/user").Subrouter()
	user_api.HandleFunc("/{id}/", Profile)

	handler := cors_middle.Handler(r)

	http.ListenAndServe(":8080", handler)
}
