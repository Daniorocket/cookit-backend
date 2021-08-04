package lib

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var jwtKey = []byte(os.Getenv("JWT_KEY"))

type Token struct {
	Token          interface{} `json:"token"`
	ExpirationTime string      `json:"expirationTime"`
}

func CreateJWT(username string, expireTime time.Time) (string, error) {
	claims := &claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func VerifyAndReturnJWT(w http.ResponseWriter, r *http.Request) (string, string, error) {

	tokenFromRequest := r.Header.Get("Authorization")
	if len(tokenFromRequest) < 10 {
		return "", "", errors.New("Failed to read token")
	}
	dividedToken := strings.Split(tokenFromRequest, " ")
	if len(dividedToken) != 2 {
		return "", "", errors.New("Failed to divide token")
	}
	tokenFromRequest = dividedToken[1]
	claims := &claims{}

	tkn, err := jwt.ParseWithClaims(tokenFromRequest, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", "", errors.New("Invalid signature of token")
		}
		return "", "", errors.New("Invalid token")
	}
	if !tkn.Valid {
		return "", "", errors.New("Invalid token")
	}
	return tokenFromRequest, claims.Username, nil
}
func RenewJWT(w http.ResponseWriter, r *http.Request, tknStr string) (string, string, error) {
	claims := &claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", "", errors.New("Invalid signature of token")
		}
		return "", "", errors.New("Invalid token")
	}
	if !tkn.Valid {
		return "", "", errors.New("Invalid token")
	}
	expirationTime := time.Now().Add(15 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", "", errors.New("Failed to generate token")
	}
	return tokenString, expirationTime.UTC().String(), nil
}
