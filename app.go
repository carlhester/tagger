package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type App struct {
	db  *myDB
	log *log.Logger
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

	entries := app.db.all()
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
