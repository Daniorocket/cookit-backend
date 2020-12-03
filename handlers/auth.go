package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Daniorocket/cookit-backend/models"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var jwtKey = []byte("my_secret_key")

func (d *Handler) Login(w http.ResponseWriter, r *http.Request) {
	login := &models.Login{}
	err := json.NewDecoder(r.Body).Decode(login)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Failed to decode body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	passDB, err := models.GetPasswordByUsernameOrEmail(d.Sess, d.DatabaseName, login.Username)
	if err != nil {
		log.Println("Failed:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passDB), []byte(login.Password)); err != nil {
		log.Println("Invalid username or password", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: login.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(Token{
		Token: tokenString,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func (d *Handler) Register(w http.ResponseWriter, r *http.Request) {
	cred := &models.Credentials{}
	err := json.NewDecoder(r.Body).Decode(cred)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Failed to decode body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cred.Password), 8)
	if err != nil {
		log.Println("Failed to hash password using bcrypt:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := &models.User{
		ID:       uuid.NewV4().String(),
		Email:    cred.Email,
		Username: cred.Username,
		Password: string(hashedPassword),
	}

	if err := models.RegisterUser(d.Sess, d.DatabaseName, user); err != nil {
		log.Println("Failed:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/v1/login", http.StatusCreated)
}
func (d *Handler) Renew(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("token").(string)
	fmt.Println("Token to renew:", tkn)

}
