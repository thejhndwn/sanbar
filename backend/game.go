package main

import (
	"net/http"
	"fmt"
	"context"
	"encoding/json"
	"time"
	"log"
)	

type NewGamePayload struct {
	Target int `json:"target"`
	NumCards int `json:"num_cards"`
	GameType string `json:"game_type"`
}


func MakeSurvival(dbm *DatabaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		log.Println("Received the game request")
		fmt.Fprint(w, "You entered the Survival")

		var payload NewGamePayload
		if err := json.NewDecoder(r.Body).Decode(&payload);
		err != nil {
			http.Error(w, "INvalid JSON",
			http.StatusBadRequest)
			return 
		}

		defer r.Body.Close()
		target := payload.Target
		numCards := payload.NumCards
		// gameType := payload.GameType

	    type contextKey string
		var UserKey contextKey = "user"

		user_id := r.Context().Value(UserKey)
		combos := GetCombos(numCards, target, dbm)
		var id string
		// todo: add numcards and target later
		err := dbm.pool.QueryRow(r.Context(),
			"INSERT INTO solo_survival_games(user_id, combos) VALUES ($1, $2) RETURNING id",
			user_id, combos,
		).Scan(&id)

		if err != nil {
			fmt.Println("make survival failed")
		}

		response := map[string]string{
			"id": id,
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if nil := json.NewEncoder(w).Encode(response); err != 
		nil {
			fmt.Println("JSON encode error:", err)
		}
	}


}

//user submits problem, if no more problems left trigger end sequence, otherwise update game data and serve next problem
// user submits solution
// get current time and calculate score from previous time
// update solve_timestmaps with the current time
// increment the timestamp

// attempt to send the next combo
// if the index is creater than len(combo) then we have done all teh problems. Transition into exit mode (tell the frontend to move into end link, or maybe just let them handle all that)
func Submit(dbm *DatabaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx:= context.Background()
		gameID:= "id"
		var gameIndex int
		var combos []string

		currentTime := time.Now()
		score := CalculateScore(currentTime)


		err := dbm.pool.QueryRow(ctx,
			`UPDATE solo_survival_games
			SET
				solve_timestamps = array_append(solve_timestamps, $1) 
				game_index = game_index +1
				score = score + $2
				scores = array_append(scores, $2)
				updated_at = NOW(),
			WHERE id = $3
			RETURNING combos, game_index
			`, currentTime, score, gameID).Scan(&combos, &gameIndex)

		if err != nil {
			fmt.Println("There was an issue with the db call to start the game and serve the first problem")
		}
			
		// start end procedure, make some handshake to end the thing. or return nothing and let the frontend move you to the end screen
		// we have to update the game state to 'completed'
		if gameIndex > len(combos) {
			return
		}

		// might have to add the header and status
		response := map[string]string{
			"combo": combos[gameIndex],
		}
		json_err:= json.NewEncoder(w).Encode(response)

		if json_err != nil {
			fmt.Println("There was an issue writing the json to the writer")
		}
	}
	
}

// game is ended, get the game stats?
func Get(dbm *DatabaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx:= context.Background()
		gameID:= "id"
		var score int

		err:= dbm.pool.QueryRow(ctx, `
			SELECT score FROM solo_survival_games
			WHERE id = "$1"
			`, gameID).Scan(&score)
		if err != nil {
			fmt.Println("There was an error getting the game stats in Get")
		}

		response := map[string]string{
			"score": fmt.Sprintf("%d", score), 
		}
		json_err:= json.NewEncoder(w).Encode(response)
		if json_err != nil {
			fmt.Println("There was an issue writing the json to the writer")
		}
	}
} 


// start button was pressed, get first problem
func Start(dbm *DatabaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		ctx:= context.Background()

		gameID:= "id"
		var gameIndex int
		var combos []string


		err := dbm.pool.QueryRow(ctx,
			`UPDATE solo_survival_games
			SET
				updated_at = NOW(),
				start_time = NOW()
			WHERE id = $1
			RETURNING combos, game_index
			`, gameID).Scan(&combos, &gameIndex)

		if err != nil {
			fmt.Println("There was an issue with the db call to start the game and serve the first problem")
		}


		// might have to add the header and status
		response := map[string]string{
			"combo": combos[gameIndex],
		}
		json_err:= json.NewEncoder(w).Encode(response)

		if json_err != nil {
			fmt.Println("There was an issue writing the json to the writer")
			return
		}
	} 
}


// user clicked end early, update game state
func End(dbm *DatabaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		ctx:= context.Background()
		gameID:= "id"
		err := dbm.pool.QueryRow(ctx,
			`UPDATE solo_survival_games
			SET
				updated_at = NOW(),
				end_time = NOW(),
				solve_timestamps = array_append(solve_timestamps, NOW()),
				scores = array_append(scores, 0),
				status = 'completed'
			WHERE id = $1
			`, gameID)

		if err != nil {
			fmt.Println("There was an error with the db call in End")
			return
		}
	} 
}

// user opts to skip the problem, update the game and serve the next problem
func Skip(dbm *DatabaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		ctx:= context.Background()
		gameID:= "id"
		var gameIndex int
		var combos []string

		currentTime := time.Now()
		score := CalculateSkipScore()


		err := dbm.pool.QueryRow(ctx,
			`UPDATE solo_survival_games
			SET
				solve_timestamps = array_append(solve_timestamps, $1) 
				game_index = game_index +1
				score = score + $2
				scores = array_append(scores, $2)
				updated_at = NOW(),
			WHERE id = $3
			RETURNING combos, game_index
			`, currentTime, score, gameID).Scan(&combos, &gameIndex)

		if err != nil {
			fmt.Println("There was an issue with the db call to start the game and serve the first problem")
		}
			
		// start end procedure, make some handshake to end the thing. or return nothing and let the frontend move you to the end screen
		// we have to update the game state to 'completed'
		if gameIndex > len(combos) {
			return
		}

		// might have to add the header and status
		response := map[string]string{
			"combo": combos[gameIndex],
		}
		json_err:= json.NewEncoder(w).Encode(response)

		if json_err != nil {
			fmt.Println("There was an issue writing the json to the writer")
		}

	} 
}

/**
func Break(dbm *DatabaseManager) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
	}
}

func Continue(dbm *DatabaseManager) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
	}
}


// TODO later, potentially useful for head to head mode
func Ready(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You entered the Ready")
}
**/
