package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
	dataField  = "data"
)

func main() {
	flag.Parse()

	// determine note data
	var err error
	var data []byte
	if *file != "" {
		data, err = ioutil.ReadFile(*file)
		fatalif(err)
	} else if flag.NArg() == 1 {
		data = []byte(flag.Arg(0))
	} else {
		log.Fatal("one and only one arg must be provided")
	}

	// open database
	db, err := sql.Open("sqlite3", *dbpath)
	fatalif(err)
	defer db.Close()

	// create table
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (%v INTEGER, %v BLOB);",
		notesTable, timeField, dataField)
	_, err = db.Exec(sql)
	fatalif(err)

	// insert note into db
	sql = fmt.Sprintf("INSERT INTO %v VALUES (?,?);", notesTable)
	_, err = db.Exec(sql, time.Now().Unix(), data)
	fatalif(err)

	fmt.Printf("Data dumped successfully to '%v'\n", dbname)
}

func fatalif(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
