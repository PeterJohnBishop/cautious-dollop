package postgresdb

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

func GenerateTokens(userID string) (accessToken, refreshToken string, err error) {
	accessClaims := UserClaims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	refreshClaims := UserClaims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessToken, err = at.SignedString([]byte(AccessSecret))
	if err != nil {
		return
	}
	refreshToken, err = rt.SignedString([]byte(RefreshSecret))
	return
}

func ValidateToken(tokenStr string, isRefresh bool) (*UserClaims, error) {
	fmt.Printf("Validating token: %s\n", tokenStr)
	secret := AccessSecret
	if isRefresh {
		fmt.Println("Using refresh token secret")
		secret = RefreshSecret
	}

	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil // ✅ convert string to []byte here
	})
	if err != nil {
		fmt.Printf("Error parsing token: %v\n", err)
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateToken(tokenStr, false)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.ID)
		c.Next()
	}
}
