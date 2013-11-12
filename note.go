package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "code.google.com/p/go-sqlite/go1/sqlite3"
)

var defaultdb = filepath.Join(os.Getenv("HOME"), "personal-notes.sqlite")

var (
	file   = flag.String("f", "", "use file contents as the note text")
	dbpath = flag.String("db", defaultdb, "path to notes database")
)

const (
	notesTable = "notes"
	timeField  = "unixtime"
	metaField  = "meta"
	dataField  = "data"
)

func main() {
	flag.Parse()

	if flag.NArg() > 1 {
		log.Fatal("only one arg must be provided")
	}

	// get file and/or meta data
	var err error
	data := []byte{}
	meta := ""
	if *file != "" {
		data, err = ioutil.ReadFile(*file)
		fatalif(err)
	}
	if flag.NArg() == 1 {
		meta = flag.Arg(0)
	}

	// open database
	db, err := sql.Open("sqlite3", *dbpath)
	fatalif(err)
	defer db.Close()

	// create table
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (%v INTEGER, %v TEXT, %v BLOB);",
		notesTable, timeField, metaField, dataField)
	_, err = db.Exec(sql)
	fatalif(err)

	// insert note into db
	sql = fmt.Sprintf("INSERT INTO %v VALUES (?,?,?);", notesTable)
	_, err = db.Exec(sql, time.Now().Unix(), meta, data)
	fatalif(err)

	fmt.Printf("Data dumped successfully to '%v'\n", *dbpath)
}

func fatalif(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
