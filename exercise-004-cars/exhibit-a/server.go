package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

type Vehicle struct {
	Name  string
	Count int
}

type Vehicles struct {
	List []Vehicle
}

type View struct {
	Username string
	Vehicles Vehicles
}

var joinT = template.Must(template.ParseFiles("templates/join.html"))
var playT = template.Must(template.ParseFiles("templates/play.html"))
var viewInstances []View

//var names []string

func home(w http.ResponseWriter, r *http.Request) {
	// display the home page
	joinT.Execute(w, nil)
}

func join(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")

	if username != "" {
		// store in a cookie
		cookie := http.Cookie{Name: "username", Value: username, Expires: inOneYear()}
		http.SetCookie(w, &cookie)

		// add to our list
		viewInstances = append(viewInstances, View{Username: username})

		// redirect browser to the play view.
		http.Redirect(w, r, "/play", http.StatusSeeOther)
	} else {
		// redirect browser back to the home view.
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func play(w http.ResponseWriter, r *http.Request) {
	username, err := r.Cookie("username")
	if err != nil {
		// redirect browser back to the home view.
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// figure out which instance of the view structure goes with this page
	for _, viewInstance := range viewInstances {
		if viewInstance.Username == username.Value {
			// display the play page on this instance
			playT.Execute(w, viewInstance)
			return
		}
	}
	// if we got here, something went wrong - redirect back to the home page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func add(w http.ResponseWriter, r *http.Request) {
	username, err := r.Cookie("username")
	if err != nil {
		// redirect browser back to the home view.
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	vehicle := r.Form.Get("vehicle")
	speed := r.Form.Get("speed")
	vehicle = vehicle + ":" + speed

	// figure out which instance of the view structure goes with this page.
	// Note that viewInstance is a copy by value, do not try to modify it in the loop.
	for i, viewInstance := range viewInstances {
		// if this instance matches our username.
		if viewInstance.Username == username.Value {
			// Loop through the existing vehicles
			for j, v := range viewInstance.Vehicles.List {
				if v.Name == vehicle {
					// Increment the count if the vehicle already exists
					viewInstances[i].Vehicles.List[j].Count++

					// redirect to play
					http.Redirect(w, r, "/play", http.StatusSeeOther)
					return
				}
			}

			// If the vehicle doesn't exist, append a new vehicle with count 1
			viewInstances[i].Vehicles.List = append(viewInstances[i].Vehicles.List, Vehicle{Name: vehicle, Count: 1})
			logVehiclesList(&viewInstances[i])

			// redirect to play
			http.Redirect(w, r, "/play", http.StatusSeeOther)
			return
		}
	}
	// if we got here, something went wrong - redirect back to the home page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func logVehiclesList(viewInstance *View) {
	log.Println("Vehicles List:")
	for _, vehicle := range viewInstance.Vehicles.List {
		log.Printf("Name: %s, Count: %d", vehicle.Name, vehicle.Count)
	}
}

func exit(w http.ResponseWriter, r *http.Request) {
	username, err := r.Cookie("username")
	if err != nil {
		// redirect browser back to the home view.
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// find that user and delete from the list
	for i, viewInstance := range viewInstances {
		if viewInstance.Username == username.Value {
			// remove this user from the list
			viewInstances = append(viewInstances[:i], viewInstances[i+1:]...)
		}
	}

	// delete the cookie
	username.MaxAge = -1
	http.SetCookie(w, username)

	// redirect browser back to the home view.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func inOneYear() time.Time {
	return time.Now().AddDate(1, 0, 0)
}

func poke(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: "username", Value: "gopher", Expires: inOneYear()}
	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, "Just set cookie named 'username' set to 'gopher'")
}

func peek(w http.ResponseWriter, r *http.Request) {
	username, err := r.Cookie("username")
	if err != nil {
		fmt.Fprintf(w, "Could not find cookie named 'username'")
		return
	}
	fmt.Fprintf(w, "You have a cookie named 'username' set to '%s'", username)
}

func hide(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: "username", MaxAge: -1}
	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, "The cookie named 'username' should be gone!")
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/poke", poke)
	http.HandleFunc("/peek", peek)
	http.HandleFunc("/hide", hide)
	http.HandleFunc("/play", play)
	http.HandleFunc("/join", join)
	http.HandleFunc("/exit", exit)
	http.HandleFunc("/add", add)

	// Serve files from the "public" directory at the "/public/" URL path
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.ListenAndServe(":8080", nil)
}
