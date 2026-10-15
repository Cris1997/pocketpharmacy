package main

import (
	"encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
	http.HandleFunc("/hola", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "Hola mundo")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
    port = "8080"
    }
	
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
