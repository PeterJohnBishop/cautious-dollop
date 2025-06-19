package servegin

import (
	"database/sql"
	"log"

	"cautious-dollop/main.go/postgresdb"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func InitServer() {

	// connecting to Postgres
	db := postgresdb.ConnectPSQL(db)
	err := db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	postgresdb.CreateUsersTable(db)

	// creating gin server
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	protected := router.Group("/api")
	protected.Use(postgresdb.JWTMiddleware())

	addOpenUserRoutes(router, db)
	addProtectedUserRoutes(protected, db)
	addProtectedFileRoutes(protected)

	router.Run(":8080")
}
