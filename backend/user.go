package main

import (
	"net/http"
	"fmt"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You entered Profile")
}
