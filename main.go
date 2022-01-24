package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	db  *myDB
	log *log.Logger
}

type myDB struct {
	db  *sql.DB
	log *log.Logger
}

type row struct {
	id   int
	link string
	tags []string
}

func (d *myDB) searchTags(query string) []Entry {
	allRows := []row{}
	q := fmt.Sprintf("SELECT * FROM tags where tags like '%%%s%%'", query)
	d.log.Printf("Query: %+v", q)
	rows, err := d.db.Query(q)
	if err != nil {
		d.log.Fatal(err)
	}

	for rows.Next() {
		var r row
		var tmpTags string
		err := rows.Scan(&r.id, &r.link, &tmpTags)
		if err != nil {
			d.log.Fatal(err)
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
	var rows = &sql.Rows{}
	if link != "" && tags != "" {
		q := fmt.Sprintf("insert into tags (link, tags) values ('%s', '%s')", link, tags)
		d.log.Printf("Insert: %+v", q)
		var err error
		rows, err = d.db.Query(q)
		if err != nil {
			d.log.Fatal(err)
		}
	} else {
		return d.all()
	}

	for rows.Next() {
		var r row
		var tmpTags string
		err := rows.Scan(&r.id, &r.link, &tmpTags)
		if err != nil {
			d.log.Fatal(err)
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
		d.log.Fatal(err)
	}

	for rows.Next() {
		var r row
		var tmpTags string
		err := rows.Scan(&r.id, &r.link, &tmpTags)
		if err != nil {
			d.log.Fatal(err)
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
	if req.URL.Path == "/data" {
		e := Entry{
			Link: "TestLink",
			Tags: []string{"test1", "test2"},
		}
		b, err := json.Marshal(e)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, string(b))
		return
	}

	entries := func(req *http.Request) []Entry {
		switch req.URL.Path {
		case "/search":
			if req.Method == http.MethodPost {
				req.ParseForm()
				app.log.Printf("ParseForm: %+v\n", req.Form)
				t := req.Form["search"]
				return app.db.searchTags(t[0])
			} else {
				q, _ := url.ParseQuery(req.URL.RawQuery)
				return app.db.searchTags(fmt.Sprintf(q["tag"][0]))
			}
		case "/add":
			req.ParseForm()
			app.log.Printf("\n%+v\n", req.Form)
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
		app.log.Fatal(err)
	}
	tmpl.Execute(w, pageData)
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
		db:  &myDB{db: db, log: logger},
		log: logger,
	}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", app.handler)
	http.HandleFunc("/add", app.handler)
	http.HandleFunc("/data", app.handler)
	http.HandleFunc("/search", app.handler)

	fmt.Println("Listening on ", url)
	log.Fatal(http.ListenAndServe(url, nil))
}
