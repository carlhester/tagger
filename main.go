package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type myDB struct {
	db *sql.DB
}

type row struct {
	id   int
	link string
	tags []string
}

func (d *myDB) all() []row {
	result := []row{}
	rows, err := d.db.Query("SELECT * FROM tags")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var r row
		var tmpTags string
		err := rows.Scan(&r.id, &r.link, &tmpTags)
		if err != nil {
			panic(err)
		}

		r.tags = strings.Fields(tmpTags)

		result = append(result, r)
	}
	return result
}

type Entry struct {
	Link string
	Tags []string
}

type PageData struct {
	PageTitle string
	Entries   []Entry
}

func hello(w http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}
	database := myDB{db: db}
	allRows := database.all()

	entries := []Entry{}
	for _, d := range allRows {
		entries = append(entries, Entry{Link: d.link, Tags: d.tags})
	}

	pageData := PageData{
		PageTitle: "Tagged Sites",
		Entries:   entries,
	}

	tmpl, err := template.ParseFiles("./layout.tmpl")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, pageData)
}

func main() {

	http.HandleFunc("/", hello)
	http.ListenAndServe(":9191", nil)
}
