package main

import (
	"net/http"
	"fmt"
	"context"
	"encoding/json"
	"time"
)	

type NewGamePayload struct {
	Target int `json:"target"`
	NumCards int `json:"num_cards"`
	GameType string `json:"game_type"`
}

type StartGamePayload struct {
	GameId string `json:"gameId"`
}

type GenericGamePayload struct {
	GameId string `json:"gameId"`
}


func MakeSurvival(dbm *DatabaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		if r.Method == "OPTIONS" { 
			return 
		}

		var payload NewGamePayload
		if err := json.NewDecoder(r.Body).Decode(&payload);
		err != nil {
			http.Error(w, "Invalid JSON",
			http.StatusBadRequest)
			return 
		}

		defer r.Body.Close()
		target := payload.Target
		numCards := payload.NumCards

		user_id := r.Context().Value(userIDKey)
		combos := GetCombos(numCards, target, dbm)
		var id string

		if user_id == nil{
			fmt.Println("There was an error getting your user_id in makeSurival")
		}
		fmt.Println("In makesurvival we got the userid:", user_id)
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

	if err:= json.NewEncoder(w).Encode(response); err != nil {
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
		fmt.Println("we are making a submission")
		
		var payload GenericGamePayload
		err:= json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			fmt.Println("There was an error parsing the payload in submitgame:", err)
		}

		var gameIndex int
		var combos []string

		currentTime := time.Now()
		score := CalculateScore(currentTime) // TODO: grab previous timestamp to actually calculate real score

		queryerr := dbm.pool.QueryRow(r.Context(),
		`UPDATE solo_survival_games
		SET
		solve_timestamps = array_append(solve_timestamps, $1),
		game_index = game_index +1,
		scores = array_append(scores, $2),
		score = score + $2,
		problem_remaining = problem_remaining - 1,
		problem_solved = problem_solved + 1,

		updated_at = NOW()
		WHERE id = $3
		RETURNING combos, game_index
		`, currentTime, score, payload.GameId).Scan(&combos, &gameIndex)

		if queryerr != nil {
			fmt.Println("There was an issue submitting the solution")
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		if gameIndex >= len(combos) {
			response := map[string]string{
				"status": "ended",
			}
			json_err:= json.NewEncoder(w).Encode(response)
			if json_err != nil {
				fmt.Println("There was an issue writing the json to the writer")
			}
			err := dbm.pool.QueryRow(r.Context(),
			`UPDATE solo_survival_games
			SET
			updated_at = NOW(),
			end_time = NOW(),
			status = 'completed'
			WHERE id = $1
			`, payload.GameId)

			if err != nil {
				fmt.Println("There was an error with the db call in End")
			}
			return 
		}

		// might have to add the header and status
		response := map[string]string{
			"status": "ongoing",
			"combo": combos[gameIndex],
		}
		json_err:= json.NewEncoder(w).Encode(response)

		if json_err != nil {
			fmt.Println("There was an issue writing the json to the writer in submit:", json_err)
		}
	}

}

// game is ended, get the game stats?
func Get(dbm *DatabaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx:= context.Background()
		gameID:= "id"
		var scores []int

		err:= dbm.pool.QueryRow(ctx, `
		SELECT scores FROM solo_survival_games
		WHERE id = "$1"
		`, gameID).Scan(&scores)
		if err != nil {
			fmt.Println("There was an error getting the game stats in Get")
		}

		score := 0
		for _, num := range scores {
			score += num
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

// TODO: need to add userID-gameID confirmation
func Start(dbm *DatabaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		fmt.Println("HitStART")

		userID := r.Context().Value(userIDKey).(string)

		var payload StartGamePayload
		err:= json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			fmt.Println("There was an error parsing the payload in gamestart")
			fmt.Println(r.Body)
			fmt.Println(payload)
		}

		var gameIndex int
		var combos []string
		fmt.Println("payload is:", payload)


		queryerr := dbm.pool.QueryRow(r.Context(),
		`UPDATE solo_survival_games
		SET
		updated_at = NOW(),
		start_time = NOW(),
		status = 'active',
		problems_remaining = array_length(combos, 1)
		WHERE id = $1 AND user_id = $2
		RETURNING combos, game_index
		`, payload.GameId, userID).Scan(&combos, &gameIndex)

		if queryerr != nil {
			fmt.Println("There was an issue with the db call to start the game and serve the first problem", queryerr)
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
		var payload GenericGamePayload
		err:= json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			fmt.Println("There was an error parsing the payload in End")
		}
		queryerr := dbm.pool.QueryRow(r.Context(),
		`UPDATE solo_survival_games
		SET
		updated_at = NOW(),
		end_time = NOW(),
		solve_timestamps = array_append(solve_timestamps, NOW()),
		scores = array_append(scores, 0),
		status = 'completed'
		WHERE id = $1
		`, payload.GameId)

		if queryerr != nil {
			fmt.Println("There was an error with the db call in End")
			return
		}
	} 
}

// user opts to skip the problem, update the game and serve the next problem
func Skip(dbm *DatabaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){

		var payload GenericGamePayload
		err:= json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			fmt.Println("There was an error parsing the payload in gamestart")
		}


		var gameIndex int
		var combos []string

		currentTime := time.Now()
		score := CalculateSkipScore()

		queryerr := dbm.pool.QueryRow(r.Context(),
		`UPDATE solo_survival_games
		SET
		solve_timestamps = array_append(solve_timestamps, $1),
		game_index = game_index +1,
		scores = array_append(scores, $2),
		updated_at = NOW()
		WHERE id = $3
		RETURNING combos, game_index
		`, currentTime, score, payload.GameId).Scan(&combos, &gameIndex)

		if queryerr != nil {
			fmt.Println("there was an error with skip db query:", queryerr)
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		if gameIndex >= len(combos) {
			response := map[string]string{
				"status": "ended",
			}
			json_err:= json.NewEncoder(w).Encode(response)
			if json_err != nil {
				fmt.Println("There was an issue writing the json to the writer in skip to end")
			}

			err := dbm.pool.QueryRow(r.Context(),
			`UPDATE solo_survival_games
			SET
			updated_at = NOW(),
			end_time = NOW(),
			status = 'completed'
			WHERE id = $1
			`, payload.GameId)

			if err != nil {
				fmt.Println("There was an error with the db call in End")
			}
			return 
		}

		// might have to add the header and status
		response := map[string]string{
			"combo": combos[gameIndex],
		}
		json_err:= json.NewEncoder(w).Encode(response)

		if json_err != nil {
			fmt.Println("There was an issue writing the json to the writer in skip to next problem")
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
