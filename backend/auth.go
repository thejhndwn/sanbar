package main	

import (
	"fmt"
	"net/http"
	"context"
	// "encoding/json"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	GuestToken string `json:"guest_token"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println( "entered loginhandler")
}

func RegisterHandler(dbm *DatabaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "inserting user into database")
		ctx := context.Background()

		/**
		var req RegisterRequest 
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			fmt.Println("invalid json in register request")
			return
		}
		**/

		// TODO add value grabbing

		// TODO add validation

		result, err := dbm.pool.Exec(ctx,
			"INSERT INTO users (username, email, password_hash, guest_token) VALUES ($1, $2, $3, $4)",
			"test_user", "test_user@mom.com", "notatestpassword", "d3554666-2372-4179-875a-051ea9c8c732",
		)
		fmt.Printf("we have made the user entry: %s", result)

		if err != nil {
			fmt.Println("register insert failed")
		}
		
	}
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
