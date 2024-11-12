package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Phone data struct
type Phone struct {
	Age     int    `json:"age"`
	ID      string `json:"id"`
	URL     string `json:"imageUrl"`
	Name    string `json:"name"`
	Snippet string `json:"snippet"`
}

var allPhones []Phone

func setup() {
	data, err := os.ReadFile("exhibit-d/phones.json")
	if err != nil {
		fmt.Println("Error reading phones.json")
		os.Exit(1)
	}

	err = json.Unmarshal(data, &allPhones)
	if err != nil {
		fmt.Println("Error in unmarshalling phones: ", err)
	}

	// pretty printing for testing
	data, err = json.MarshalIndent(&allPhones, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}

func phones(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(allPhones)
}

func main() {
	setup()
	http.HandleFunc("/phones", phones)
	http.ListenAndServe(":8080", nil)
}
