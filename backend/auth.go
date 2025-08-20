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

type contextKey string
const userIDKey contextKey = "userID"


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

func NewUser(dbm *DatabaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		fmt.Println("We are going to add this new user")
		token := r.Header.Get("Authorization")
		fmt.Println("got token:", token)

		if err:= CreateUserWithToken(token, dbm); err != nil {
			fmt.Printf("error create the user with token in newuser: %v", err)	
		}

		response := map[string]bool{
			"success": true,
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			fmt.Println("JSON encode error:", err)
		}
	}
}

func AuthUser(next http.HandlerFunc, dbm *DatabaseManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		fmt.Println("In the auth middleware")
		token := r.Header.Get("Authorization")
		fmt.Println("we found authheader:", token)
		userID := GetUserIDFromToken(token, dbm)
		fmt.Println("we found the userID:", userID)
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w,r.WithContext(ctx))
	}
}

func GetUserIDFromToken(token string, dbm *DatabaseManager) string {
	ctx:=context.TODO()
	var id string
	err := dbm.pool.QueryRow(ctx, `SELECT id FROM users WHERE guest_token=$1`, token).Scan(&id)
	if err != nil {
		fmt.Printf("error getting id from token: %s", err)
	}

	return id
}

func CreateUserWithToken(token string, dbm *DatabaseManager) error {
	ctx := context.TODO()
	_, err := dbm.pool.Exec(ctx, `INSERT INTO users (guest_token) VALUES ($1);`, token)
	fmt.Println("added user to the db")

	if err != nil {
		fmt.Println("There was an error creating the user from the token, CreateUserWithToken")
		return fmt.Errorf("we failed to create the user with the token")
	}
	return nil

}
	
