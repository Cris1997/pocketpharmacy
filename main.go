package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Producto struct {
	ID         int64   `json:"id"`
	Nombre     string  `json:"nombre"`
	Precio     float64 `json:"precio"`
	Existencia int     `json:"existencia"`
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("La variable DATABASE_URL no está configurada")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Error al configurar PostgreSQL: ", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal("No fue posible conectar con PostgreSQL: ", err)
	}

	log.Println("Conexión a PostgreSQL exitosa")

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

	http.HandleFunc("/api/productos", obtenerProductos(db))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API iniciada en el puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func obtenerProductos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.QueryContext(r.Context(),
			`SELECT id, name, price, stock
			 FROM products
			 ORDER BY id`)
		if err != nil {
			log.Printf("Error consultando productos: %v", err)
			http.Error(w, "Error al consultar productos",
				http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		productos := make([]Producto, 0)

		for rows.Next() {
			var producto Producto

			if err := rows.Scan(
				&producto.ID,
				&producto.Nombre,
				&producto.Precio,
				&producto.Existencia,
			); err != nil {
				log.Printf("Error leyendo producto: %v", err)
				http.Error(w, "Error al leer productos",
					http.StatusInternalServerError)
				return
			}

			productos = append(productos, producto)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Error recorriendo productos: %v", err)
			http.Error(w, "Error al leer productos",
				http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(productos); err != nil {
			log.Printf("Error generando JSON: %v", err)
		}
	}
}
