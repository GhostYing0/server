package util

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret []byte

type Claims struct {
	ID       string
	Username string
	Role     int
	jwt.StandardClaims
}

func GenerateToken(userid string, username string, role int) (string, error) {
	expireTime := time.Now().Add(1 * time.Hour)
	claims := Claims{
		ID:       userid,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (claims *Claims, err error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	ok := false
	if tokenClaims != nil {
		claims, ok = tokenClaims.Claims.(*Claims)
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return claims, err
}
