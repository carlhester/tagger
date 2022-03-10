package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type api struct {
	log *log.Logger
	db  *myDB
}

func (a *api) handler(w http.ResponseWriter, req *http.Request) {
	a.log.Printf("%s %s", req.Method, req.URL.Path)
	switch req.URL.Path {
	case "/data":
		a.handleData(w, req)
	case "/add":
		a.handleAdd(w, req)
	}

}

func (a *api) handleData(w http.ResponseWriter, req *http.Request) {
	a.log.Println("handleData entered")
	// e := Entry{
	// 	Link: "TestLink",
	// 	Tags: []string{"test1", "test2"},
	// }
	e := a.db.all()
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(b))
}

func (a *api) handleAdd(w http.ResponseWriter, req *http.Request) {
	var e entry
	err := json.NewDecoder(req.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v\n", e)
}

type entry struct {
	Link string
	Tags []string
}
