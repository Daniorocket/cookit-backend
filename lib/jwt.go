package lib

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var jwtKey = []byte("my_secret_key")

type Token struct {
	Token interface{} `json:"token"`
}

func CreateJWT(username string, expireTime time.Time) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return "", err
	}
	return tokenString, nil
}
func VerifyAndReturnJWT(w http.ResponseWriter, r *http.Request) (string, error) {

	tokenFromRequest := r.Header.Get("Authorization")
	if len(tokenFromRequest) < 10 {
		w.WriteHeader(http.StatusUnauthorized)
		return "", errors.New("Failed to read token")
	}
	dividedToken := strings.Split(tokenFromRequest, " ")
	if len(dividedToken) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return "", errors.New("Failed to divide token")
	}
	tokenFromRequest = dividedToken[1]
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tokenFromRequest, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return "", errors.New("Invalid signature of token")
		}
		w.WriteHeader(http.StatusBadRequest)
		return "", errors.New("Invalid token")
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return "", errors.New("Invalid token")
	}
	return tokenFromRequest, nil
}
func RenewJWT(w http.ResponseWriter, r *http.Request, tknStr string) (string, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return "", errors.New("Invalid signature of token")
		}
		w.WriteHeader(http.StatusBadRequest)
		return "", errors.New("Invalid token")
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return "", errors.New("Invalid token")
	}
	// (END) The code up-till this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	// if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return errors.New("Failed to generate token")
	// }

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "", errors.New("Failed to generate token")
	}
	return tokenString, nil
}
