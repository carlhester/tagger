package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type myDB struct {
	db  *sql.DB
	log *log.Logger
}

// func (d *myDB) searchTags(query string) []Entry {
// 	allRows := []row{}
// 	q := fmt.Sprintf("SELECT * FROM tags where tags like '%%%s%%'", query)
// 	d.log.Printf("Query: %+v", q)
// 	rows, err := d.db.Query(q)
// 	if err != nil {
// 		d.log.Fatal(err)
// 	}

// 	for rows.Next() {
// 		var r row
// 		var tmpTags string
// 		err := rows.Scan(&r.id, &r.link, &tmpTags)
// 		if err != nil {
// 			d.log.Fatal(err)
// 		}

// 		r.tags = strings.Fields(tmpTags)

// 		allRows = append(allRows, r)
// 	}
// 	entries := []Entry{}
// 	for _, d := range allRows {
// 		entries = append(entries, Entry{Link: d.link, Tags: d.tags})
// 	}
// 	return entries
// }

// func (d *myDB) add(link, tags string) []Entry {
// 	allRows := []row{}
// 	var rows = &sql.Rows{}
// 	if link != "" && tags != "" {
// 		q := fmt.Sprintf("insert into tags (link, tags) values ('%s', '%s')", link, tags)
// 		d.log.Printf("Insert: %+v", q)
// 		var err error
// 		rows, err = d.db.Query(q)
// 		if err != nil {
// 			d.log.Fatal(err)
// 		}
// 	} else {
// 		return d.all()
// 	}

// 	for rows.Next() {
// 		var r row
// 		var tmpTags string
// 		err := rows.Scan(&r.id, &r.link, &tmpTags)
// 		if err != nil {
// 			d.log.Fatal(err)
// 		}

// 		r.tags = strings.Fields(tmpTags)

// 		allRows = append(allRows, r)
// 	}
// 	entries := []Entry{}
// 	for _, d := range allRows {
// 		entries = append(entries, Entry{Link: d.link, Tags: d.tags})
// 	}
// 	return entries
// }

func (d *myDB) all() []Entry {
	fmt.Println("all")
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
	fmt.Println(len(entries))
	return entries
}
