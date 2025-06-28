package servegin

import (
	"database/sql"

	"cautious-dollop/main.go/postgresdb"

	"github.com/gin-gonic/gin"
)

func addOpenUserRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from the Go server!",
		})
	})
	r.POST("/login", func(c *gin.Context) {
		postgresdb.Login(db, c)
	})
	r.POST("/register", func(c *gin.Context) {
		postgresdb.RegisterUser(db, c)
	})
	r.GET("/refresh", func(c *gin.Context) {
		postgresdb.Refresh(c)
	})
	r.GET("/onetime/download", func(c *gin.Context) {
		handleOneTimeFileDownload(c)
	})
}

func addProtectedUserRoutes(r *gin.RouterGroup, db *sql.DB) {

	r.GET("/users", func(c *gin.Context) {
		postgresdb.GetUsers(db, c)
	})
	r.GET("/users/:id", func(c *gin.Context) {
		postgresdb.GetUserByID(db, c)
	})
	r.PUT("/users", func(c *gin.Context) {
		postgresdb.UpdateUser(db, c)
	})
	r.DELETE("/users/:id", func(c *gin.Context) {
		postgresdb.DeleteUserByID(db, c)
	})
}

func addProtectedFileRoutes(r *gin.RouterGroup) {
	r.POST("/upload", handleFileUpload)
	r.POST("/onetime/upload", func(c *gin.Context) {
		handleOneTimeFileUpload(c)
	})
	r.GET("/files", handleListFiles)
	r.GET("/download/:filename", handleFileDownload)
	r.DELETE("/delete/:filename", handleFileDelete)
}
