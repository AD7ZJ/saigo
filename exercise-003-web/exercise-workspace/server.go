package main

import (
	"html/template"
	"net/http"
)

var homeT = template.Must(template.ParseFiles("exercise-workspace/home.html"))
var names []string

func home(w http.ResponseWriter, r *http.Request) {
	// The names array of strings is passed into the template and magically rendered in the {{range.}} loop
	homeT.Execute(w, names)
}

func signup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")

	names = append(names, username)

	// after accepting the POST, redirect browser back to the home page.
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/home", home)
	http.HandleFunc("/signup", signup)
	http.ListenAndServe(":8080", nil)
}
