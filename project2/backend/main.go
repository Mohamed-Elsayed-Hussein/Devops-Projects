package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// connect opens a connection to the MySQL database using the password
// read from the Docker secret at /run/secrets/db-password.
// The returned *sql.DB is a pooled handle; callers should Close() it when done.
func connect() (*sql.DB, error) {
	bin, err := ioutil.ReadFile("/run/secrets/db-password")
	if err != nil {
		return nil, err
	}
	return sql.Open("mysql", fmt.Sprintf("root:%s@tcp(db:3306)/example", string(bin)))
}

// blogHandler is the HTTP handler for the root path.
// It queries the blog table for all titles and responds with a JSON array of titles.
func blogHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connect()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT title FROM blog")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	var titles []string
	for rows.Next() {
		var title string
		err = rows.Scan(&title)
		titles = append(titles, title)
	}
	// Encode the slice of titles as JSON and write it to the response.
	json.NewEncoder(w).Encode(titles)
}

// main is the application entrypoint.
// It prepares the database, sets up the HTTP router, and starts the server on :8000
// with request logging enabled.
func main() {
	log.Print("Prepare db...")
	if err := prepare(); err != nil {
		log.Fatal(err)
	}

	log.Print("Listening 8000")
	r := mux.NewRouter()
	r.HandleFunc("/", blogHandler)
	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r)))
}

// prepare ensures the database is ready for use by:
// 1) waiting for the DB to become reachable,
// 2) recreating the blog table, and
// 3) seeding it with a few sample blog posts.
func prepare() error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()

	// Wait up to ~60 seconds for the database to be reachable (useful in containerized setups).
	for i := 0; i < 60; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	// Recreate the blog table for a clean state on each start.
	if _, err := db.Exec("DROP TABLE IF EXISTS blog"); err != nil {
		return err
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS blog (id int NOT NULL AUTO_INCREMENT, title varchar(255), PRIMARY KEY (id))"); err != nil {
		return err
	}

	// Seed the table with a few sample posts.
	for i := 0; i < 5; i++ {
		if _, err := db.Exec("INSERT INTO blog (title) VALUES (?);", fmt.Sprintf("Blog post #%d", i)); err != nil {
			return err
		}
	}
	return nil
}