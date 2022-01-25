package main

import (
	"html/template"
	"log"
	"net/http"
)

type App struct {
	db  *myDB
	log *log.Logger
}

func (app *App) handler(w http.ResponseWriter, req *http.Request) {
	pageData := PageData{
		PageTitle: "Tagged Sites",
	}

	tmpl, err := template.ParseFiles("./layout.tmpl")
	if err != nil {
		app.log.Fatal(err)
	}
	tmpl.Execute(w, pageData)
}
