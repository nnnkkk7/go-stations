package main

import (
	"log"
	"net/http"
	"os"
	"time"

	json "encoding/json"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
	"github.com/TechBowl-japan/go-stations/model"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// set http handlers
	mux := router.NewRouter(todoDB)

	// TODO: ここから実装を行う
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		en := json.NewEncoder(w)
		m := model.HealthzResponse{Message: "OK"}
		if err = en.Encode(m); err != nil {
			log.Println(err)
		}
	})
	if err = http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}

	return nil
}
