package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
var conn *pgx.Conn

type book struct {
	Author string
	Name   string
}

type books []book

func getByAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	author := vars["autor"]
	var b book
	var bk books

	rows, err := conn.Query(context.Background(), "select name, author from books where author=$1", author)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return
	}
	for rows.Next() {
		err := rows.Scan(&b.Name, &b.Author)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			return
		}
		bk = append(bk, b)
	}
	fmt.Println(bk)
	json.NewEncoder(w).Encode(bk)
}

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	router := mux.NewRouter()

	router.HandleFunc("/api/get/{autor}", getByAuthor)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8800",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
