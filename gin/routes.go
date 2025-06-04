package servegin

import (
	"database/sql"

	"cautious-dollop/main.go/pgdb"

	"github.com/gin-gonic/gin"
)

func addOpenUserRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/login", func(c *gin.Context) {
		pgdb.Login(db, c)
	})
	r.POST("/register", func(c *gin.Context) {
		pgdb.RegisterUser(db, c)
	})
	r.GET("/refresh", func(c *gin.Context) {
		pgdb.Refresh(c)
	})
}

func addProtectedUserRoutes(r *gin.RouterGroup, db *sql.DB) {

	r.GET("/users", func(c *gin.Context) {
		pgdb.GetUsers(db, c)
	})
	r.GET("/users/:id", func(c *gin.Context) {
		pgdb.GetUserByID(db, c)
	})
	r.PUT("/users", func(c *gin.Context) {
		pgdb.UpdateUser(db, c)
	})
	r.DELETE("/users/:id", func(c *gin.Context) {
		pgdb.DeleteUserByID(db, c)
	})
}

func addProtectedFileRoutes(r *gin.RouterGroup) {
	r.POST("/upload", handleFileUpload)
	r.GET("/files", handleListFiles)
	r.GET("/download/:filename", handleFileDownload)
	r.DELETE("/delete/:filename", handleFileDelete)
}
