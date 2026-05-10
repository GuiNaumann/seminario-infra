package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"biblioteca-api/handlers"
	"biblioteca-api/repository"

	_ "github.com/lib/pq"
)

// ─── Setup ───────────────────────────────────────────────────────────────────

func main() {
	db := connectDB()
	defer db.Close()

	repo := repository.NewBookRepository(db)
	bh := handlers.NewBookHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/books", bh.BooksHandler)
	mux.HandleFunc("/books/", bh.BookHandler)

	port := getenv("PORT", "8080")
	log.Printf("Servidor rodando em :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func connectDB() *sql.DB {
	dsn := getenv("DATABASE_URL",
		"host=localhost port=5432 user=postgres password=postgres dbname=biblioteca sslmode=disable")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("erro ao abrir conexão:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("banco indisponível:", err)
	}
	log.Println("Conectado ao PostgreSQL")
	return db
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
