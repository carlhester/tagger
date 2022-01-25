package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type row struct {
	id   int
	link string
	tags []string
}

type Entry struct {
	Link string
	Tags []string
}

type PageData struct {
	PageTitle string
	Entries   []Entry
}

func main() {

	host := "0.0.0.0"
	port := "9191"
	url := fmt.Sprintf("%s:%s", host, port)

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(os.Stdout, "", log.Ltime)

	app := &App{
		log: logger,
		db: &myDB{
			db:  db,
			log: logger,
		},
	}

	api := api{log: logger}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", app.handler)
	http.HandleFunc("/add", api.handler)
	http.HandleFunc("/data", api.handler)

	fmt.Println("Listening on ", url)
	log.Fatal(http.ListenAndServe(url, nil))
}
