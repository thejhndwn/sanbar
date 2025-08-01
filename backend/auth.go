package main	

import (
	"fmt"
	"net/http"
	"context"
	"encoding/json"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	GuestToken string `json:"guest_token"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}


func LoginHandler( dbm *DatabaseManager) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		fmt.Println( "entered loginhandler")
		ctx:= context.Background()

		var req LoginRequest
		err  := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			fmt.Println("There was an error in dejsoning the login request:", req)
		}

		result, err:= dbm.pool.Exec(ctx, 
			`SELECT (username, password) 
			FROM users
			WHERE username=$1 AND password=$2
			`, req.Username, req.Password)
		if err != nil {
			fmt.Println("There was an error in database logging in:", err)
		}

		fmt.Println("we got this result:", result)
	}
}

func RegisterHandler(dbm *DatabaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "inserting user into database")
		ctx := context.Background()

		var req RegisterRequest 
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			fmt.Println("invalid json in register request")
			return
		}

		result, err:= dbm.pool.Exec(ctx,
			`UPDATE users
			SET
				updated_at = NOW(),
				is_guest = false,
				last_active= NOW(),
			WHERE id = $1
			RETURNING combos, game_index
			`)
		if err != nil {
			fmt.Println("There is an error updating the register info:", err)
		}

		fmt.Println("Got this response from register sql", result)

		
	}
}

func AuthorizeUser() {

}

func GetUserIDFromToken(token string, dbm *DatabaseManager) string {
	ctx:=context.Background()
	var id string
	err := dbm.pool.QueryRow(ctx, `SELECT id FROM users WHERE guest_token="$1"`, token).Scan(&id)
	if err != nil {
		fmt.Printf("error getting id from token: %s", err)
	}

	return id
}
