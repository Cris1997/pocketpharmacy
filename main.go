package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Producto struct {
	ID         int     `json:"id"`
	Nombre     string  `json:"nombre"`
	Precio     float64 `json:"precio"`
	Existencia int     `json:"existencia"`
}

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

	http.HandleFunc("/api/productos", obtenerProductos)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func obtenerProductos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	productos := []Producto{
		{ID: 1, Nombre: "Paracetamol 500 mg", Precio: 35.50, Existencia: 20},
		{ID: 2, Nombre: "Alcohol 250 ml", Precio: 28.00, Existencia: 15},
		{ID: 3, Nombre: "Gasas estériles", Precio: 42.00, Existencia: 10},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(productos); err != nil {
		log.Printf("Error al generar JSON: %v", err)
	}
}
