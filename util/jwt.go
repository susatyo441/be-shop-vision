package util

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTClaims adalah struktur payload dalam token JWT
type JWTClaims struct {
	Session string `json:"session"`
	ID      string `json:"id"`
	StoreID string `json:"store"`
	jwt.StandardClaims
}

// GenerateJWT membuat token JWT dengan `storeId`, `userId`, dan `session`
func GenerateJWT(userID string, storeID string, withExp bool) (string, error) {
	var jwtKey = []byte(os.Getenv("JWT_KEY"))
	fmt.Println(jwtKey)
	claims := JWTClaims{
		Session: userID,
		ID:      userID,
		StoreID: storeID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(), // Tambahkan `iat`
		},
	}

	if withExp {
		claims.StandardClaims.ExpiresAt = time.Now().Add(24 * time.Hour).Unix() // Tambahkan `exp` jika `withExp == true`
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
