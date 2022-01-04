package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	db *myDB
}

type myDB struct {
	db *sql.DB
}

type row struct {
	id   int
	link string
	tags []string
}

func (d *myDB) searchTags(query string) []Entry {
	allRows := []row{}
	q := fmt.Sprintf("SELECT * FROM tags where tags like '%%%s%%'", query)
	fmt.Printf("%+v", q)
	rows, err := d.db.Query(q)
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

		allRows = append(allRows, r)
	}
	entries := []Entry{}
	for _, d := range allRows {
		entries = append(entries, Entry{Link: d.link, Tags: d.tags})
	}
	return entries
}

func (d *myDB) add(link, tags string) []Entry {
	allRows := []row{}
	q := fmt.Sprintf("insert into tags (link, tags) values ('%s', '%s')", link, tags)
	fmt.Printf("%+v", q)
	rows, err := d.db.Query(q)
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

		allRows = append(allRows, r)
	}
	entries := []Entry{}
	for _, d := range allRows {
		entries = append(entries, Entry{Link: d.link, Tags: d.tags})
	}
	return entries
}

func (d *myDB) all() []Entry {
	allRows := []row{}
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

		allRows = append(allRows, r)
	}

	entries := []Entry{}
	for _, d := range allRows {
		entries = append(entries, Entry{Link: d.link, Tags: d.tags})
	}
	return entries
}

type Entry struct {
	Link string
	Tags []string
}

type PageData struct {
	PageTitle string
	Entries   []Entry
}

func (app *App) handler(w http.ResponseWriter, req *http.Request) {

	entries := func(req *http.Request) []Entry {
		switch req.URL.Path {
		case "/search":
			req.ParseForm()
			fmt.Printf("\n%+v\n", req.Form)
			t := req.Form["search"]
			return app.db.searchTags(t[0])
		case "/add":
			req.ParseForm()
			fmt.Printf("\n%+v\n", req.Form)
			l := req.Form["link"]
			t := req.Form["tags"]
			return app.db.add(l[0], t[0])
		case "/":
			return app.db.all()
		}
		return app.db.all()
	}(req)

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
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}

	app := &App{
		db: &myDB{db: db},
	}

	http.HandleFunc("/", app.handler)
	http.HandleFunc("/add", app.handler)
	http.HandleFunc("/search", app.handler)
	fmt.Println("Listening on 0.0.0.0:9191")
	http.ListenAndServe("0.0.0.0:9191", nil)
}
