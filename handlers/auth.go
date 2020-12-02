package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Daniorocket/cookit-backend/models"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (d *Handler) Login(w http.ResponseWriter, r *http.Request) {

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
	}
}
func (d *Handler) Verify(w http.ResponseWriter, r *http.Request) {

}
