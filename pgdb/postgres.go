package pgdb

import (
	"database/sql"
	"fmt"
	"log"
)

func ConnectPSQL(db *sql.DB) *sql.DB {

	initEnv()

	host := postgresHost
	port := postgresPort
	user := postgresUser
	password := postgresPassword
	dbname := postgresDBName
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// fmt.Println("Connecting with:", psqlInfo)

	mydb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = mydb.Ping()
	if err != nil {
		panic(err)
	}

	log.Printf("[CONNECTED] to Postgres on :%s", port)
	return mydb
}
