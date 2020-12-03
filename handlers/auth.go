package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Daniorocket/cookit-backend/lib"
	"github.com/Daniorocket/cookit-backend/models"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time

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
	tokenString, err := lib.CreateJWT(login.Username, expirationTime)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lib.Token{
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
	newTokenString, err := lib.RenewJWT(w, r, tkn)
	if err != nil {
		fmt.Println("Failed to renew token:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lib.Token{
		Token: newTokenString,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
