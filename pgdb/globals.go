package pgdb

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var postgresUser string
var postgresPassword string
var postgresDBName string
var postgresHost string
var postgresPort string
var accessSecret string
var refreshSecret string

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	getPostgresEnvs()
	getAuthEnvs()

}

func getPostgresEnvs() {

	postgresPassword = os.Getenv("PSQL_PASSWORD")
	if postgresPassword == "" {
		log.Fatal("PSQL_PASSWORD is not set in .env file")
	}
	postgresUser = os.Getenv("PSQL_USER")
	if postgresUser == "" {
		log.Fatal("PSQL_USER is not set in .env file")
	}
	postgresDBName = os.Getenv("PSQL_DBNAME")
	if postgresDBName == "" {
		log.Fatal("PSQL_DBNAME is not set in .env file")
	}
	postgresHost = os.Getenv("PSQL_HOST")
	if postgresHost == "" {
		log.Fatal("PSQL_HOST is not set in .env file")
	}
	postgresPort = os.Getenv("PSQL_PORT")
	if postgresPort == "" {
		log.Fatal("PSQL_PORT is not set in .env file")
	}

	log.Println("Postgres Environment Variables Loaded")

}

var OTDLSecret string

func getAuthEnvs() {
	accessSecret = os.Getenv("TOKEN_SECRET")
	if accessSecret == "" {
		log.Fatal("TOKEN_SECRET is not set in .env file")
	}
	refreshSecret = os.Getenv("REFRESH_TOKEN_SECRET")
	if refreshSecret == "" {
		log.Fatal("REFRESH_TOKEN_SECRET is not set in .env file")
	}
	OTDLSecret = os.Getenv("OTDL_TOKEN_SECRET")
	if OTDLSecret == "" {
		log.Fatal("OTDL_TOKEN_SECRET is not set in .env file")
	}
	log.Println("Authentication Environment Variables Loaded")
}
