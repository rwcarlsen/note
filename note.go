package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "code.google.com/p/go-sqlite/go1/sqlite3"
)

var defaultdb = filepath.Join(os.Getenv("HOME"), "personal-notes.sqlite")

var (
	query  = flag.String("q", "", "run single arg as sql query and print output")
	delim  = flag.String("delim", "", "row separator when using the query flag")
	sep    = flag.String("sep", "\t", "row separator when using the query flag")
	file   = flag.String("f", "", "use file contents as the note text")
	dbpath = flag.String("db", getdb(), "path to notes database")
)

const dbEnv = "NOTE_DATABASE"

func getdb() string {
	if os.Getenv(dbEnv) != "" {
		return os.Getenv(dbEnv)
	}
	return defaultdb
}

const (
	table     = "rawdata"
	timeField = "unixtime"
	metaField = "meta"
	dataField = "data"
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	// open database
	db, err := sql.Open("sqlite3", *dbpath)
	fatalif(err)
	defer db.Close()

	// run query
	if *query != "" {
		rows, err := db.Query(*query)
		fatalif(err)
		cols, _ := rows.Columns()
		rw := make([]interface{}, len(cols))
		for i := range rw {
			tmp := ""
			rw[i] = &tmp
		}
		for rows.Next() {
			err := rows.Scan(rw...)
			fatalif(err)
			row := make([]string, len(rw))
			for i := range rw {
				row[i] = *(rw[i]).(*string)
			}
			fmt.Print(strings.Join(row, *sep), "\n", *delim)
			if *delim != "" {
				fmt.Print("\n")
			}
		}
		fatalif(rows.Err())
		return
	}

	// get file and/or meta data
	if flag.NArg() > 2 || flag.NArg() == 0 {
		log.Fatal("1 or 2 args required: meta [note] | sql")
	}
	data := []byte{}
	meta := flag.Arg(0)
	if *file != "" {
		if flag.NArg() != 1 {
			log.Fatal("need exactly 1 meta-arg for file")
		}
		data, err = ioutil.ReadFile(*file)
		fatalif(err)
	} else {
		if flag.NArg() == 1 {
			meta = "note"
			data = []byte(flag.Arg(0))
		} else {
			data = []byte(flag.Arg(1))
		}
	}

	// create table
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (%v INTEGER, %v TEXT, %v BLOB);",
		table, timeField, metaField, dataField)
	_, err = db.Exec(sql)
	fatalif(err)

	// insert note into db
	sql = fmt.Sprintf("INSERT INTO %v VALUES (?,?,?);", table)
	_, err = db.Exec(sql, time.Now().Unix(), meta, data)
	fatalif(err)

	fmt.Printf("Data dumped successfully to '%v'\n", *dbpath)
}

func fatalif(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
