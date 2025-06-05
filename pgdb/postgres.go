package pgdb

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func ConnectPSQL(db *sql.DB) *sql.DB {
	initEnv()

	host := postgresHost
	port := postgresPort
	user := postgresUser
	password := postgresPassword
	dbname := postgresDBName

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var mydb *sql.DB
	var err error
	maxAttempts := 10

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		mydb, err = sql.Open("postgres", psqlInfo)
		if err == nil {
			err = mydb.Ping()
		}

		if err == nil {
			log.Printf("[CONNECTED] to Postgres on %s:%s", host, port)
			return mydb
		}

		log.Printf("[RETRY %d/%d] Could not connect to Postgres: %v", attempt, maxAttempts, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("Failed to connect to Postgres after %d attempts: %v", maxAttempts, err)
	return nil
}
