package main	

import (
	"fmt"
	"net/http"
)

type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    GuestID  string `json:"guest_id"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

/**
func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, welcome to login auth")
	var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    user, exists := usersDB["byemail:"+req.Email] // fake lookup
    if !exists {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    hashed, _ := passwordsDB[req.Email]
    if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(req.Password)); err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    // Generate JWT (same as register)
    expirationTime := time.Now().Add(7 * 24 * time.Hour)
    claims := &Claims{
        UserID: user.ID,
        Email:  user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenStr, err := token.SignedString(jwtSecret)
    if err != nil {
        http.Error(w, "Could not generate token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "token":   tokenStr,
        "user_id": user.ID,
    })
}

func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, welcome to register auth")
	var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Validate email/password
    if req.Email == "" || req.Password == "" {
        http.Error(w, "Email and password required", http.StatusBadRequest)
        return
    }

    if _, exists := passwordsDB[req.Email]; exists {
        http.Error(w, "Email already registered", http.StatusConflict)
        return
    }

    // Hash password
    hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    // Create user
    userID := "usr_" + generateID()
    usersDB[userID] = User{ID: userID, Email: req.Email}
    passwordsDB[req.Email] = string(hashed)

    // Migrate guest games (if guest_id provided)
    if req.GuestID != "" {
        for i := range gamesDB {
            game := gamesDB[i]
            if game["guest_id"] == req.GuestID && game["user_id"] == nil {
                game["user_id"] = userID
                game["email"] = req.Email
            }
        }
    }

    // Generate JWT
    expirationTime := time.Now().Add(7 * 24 * time.Hour)
    claims := &Claims{
        UserID: userID,
        Email:  req.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenStr, err := token.SignedString(jwtSecret)
    if err != nil {
        http.Error(w, "Could not generate token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "token":   tokenStr,
        "user_id": userID,
    })
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID")
    email := r.Context().Value("email")

    if userID == nil {
        http.Error(w, "Not authenticated", http.StatusUnauthorized)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "valid":    true,
        "user_id":  userID,
        "email":    email,
    })
}
**/

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println( "entered loginhandler")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println( "entered register")
}
