package main

import (
	"net/http"
	"fmt"
)	
func MakeSurvival(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You entered the Survival")
}


func Submit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You entered the SubmitHandler")
}
func Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You entered the Get")
}
func Ready(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You entered the Ready")
}
func Start(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You entered the Start")
}
func End(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You entered the End")
}
