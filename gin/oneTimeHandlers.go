package servegin

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"cautious-dollop/main.go/pgdb"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type FileEntry struct {
	Path      string
	ExpiresAt time.Time
}

var (
	fileStore  = make(map[string]FileEntry)
	storeMutex = sync.Mutex{}
)

func generateFileID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}

func generateDownloadJWT(fileID string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"file_id": fileID,
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(pgdb.OTDLSecret))
}

func handleOneTimeFileUpload(c *gin.Context) {
	uploadPath := "/data/uploads"

	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	freeSpace, err := getFreeSpace(uploadPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Statfs failed: %v", err)})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	if uint64(file.Size) > freeSpace {
		c.JSON(http.StatusInsufficientStorage, gin.H{"error": "Not enough disk space to save file"})
		return
	}

	expireInStr := c.DefaultPostForm("expire_in", "15m")
	expireDuration, err := time.ParseDuration(expireInStr)
	if err != nil || expireDuration <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expire_in value"})
		return
	}

	filename := filepath.Base(file.Filename)
	fileID := generateFileID()
	dst := filepath.Join(uploadPath, fmt.Sprintf("%s_%s", fileID, filename))

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Store file entry with expiration
	// Use a mutex to ensure thread-safe access to the fileStore
	// This is important if multiple uploads can happen concurrently
	// or if the store is accessed from multiple goroutines.
	storeMutex.Lock()
	fileStore[fileID] = FileEntry{
		Path:      dst,
		ExpiresAt: time.Now().Add(expireDuration),
	}
	storeMutex.Unlock()

	token, err := generateDownloadJWT(fileID, expireDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	downloadURL := fmt.Sprintf("/onetime/download?token=%s", url.QueryEscape(token))

	c.JSON(http.StatusOK, gin.H{
		"message":               "File uploaded successfully",
		"filename":              filename,
		"path":                  dst,
		"one_time_download_url": downloadURL,
		"free_space":            fmt.Sprintf("%.2f MB", float64(freeSpace)/(1024*1024)),
		"file_size":             fmt.Sprintf("%.2f MB", float64(file.Size)/(1024*1024)),
	})
}

func handleOneTimeFileDownload(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token"})
		return
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(pgdb.OTDLSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	fileID, ok := claims["file_id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing file_id"})
		return
	}

	storeMutex.Lock()
	entry, exists := fileStore[fileID]
	storeMutex.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found or already downloaded"})
		return
	}

	if time.Now().After(entry.ExpiresAt) {
		os.Remove(entry.Path)
		storeMutex.Lock()
		delete(fileStore, fileID)
		storeMutex.Unlock()
		c.JSON(http.StatusGone, gin.H{"error": "File expired"})
		return
	}

	c.File(entry.Path)

	os.Remove(entry.Path)
	storeMutex.Lock()
	delete(fileStore, fileID)
	storeMutex.Unlock()
}
