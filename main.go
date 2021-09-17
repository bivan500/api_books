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

var conn *pgx.Conn

func getByAuthor(w http.ResponseWriter, r *http.Request) {
	type book struct {
		Author string
		Name   string
	}

	type books []book

	rows, err := conn.Query(context.Background(), "select name, author from books where author=$1", mux.Vars(r)["autor"])
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return
	}

	var bookList books
	for rows.Next() {
		var bookItem book
		err := rows.Scan(&bookItem.Name, &bookItem.Author)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			return
		}
		bookList = append(bookList, bookItem)
	}
	fmt.Println(bookList)
	json.NewEncoder(w).Encode(bookList)
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
	fmt.Println("Starting server...")
	log.Fatal(srv.ListenAndServe())
}
