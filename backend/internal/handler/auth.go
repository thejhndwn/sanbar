package auth

import (
	"fmt"
	"net/http"

)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, welcome to login auth")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, welcome to register auth")
}

