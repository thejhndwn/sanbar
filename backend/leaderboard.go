package main


import (
	"net/http"
	"fmt"
)

func Survival(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "entered survival")
}
