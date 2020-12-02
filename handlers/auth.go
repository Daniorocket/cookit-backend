package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Daniorocket/cookit-backend/models"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

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
	fmt.Println("Zalogowano poprawnie!!")
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
func (d *Handler) Verify(w http.ResponseWriter, r *http.Request) {

}
