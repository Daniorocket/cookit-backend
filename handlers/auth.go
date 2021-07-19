package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Daniorocket/cookit-backend/lib"
	"github.com/Daniorocket/cookit-backend/models"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (d *Handler) Login(w http.ResponseWriter, r *http.Request) {
	login := models.Login{}
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to decode json")
		return
	}
	passDB, err := models.GetPasswordByUsernameOrEmail(d.Client, d.DatabaseName, login.Username)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Invalid username or password")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passDB), []byte(login.Password)); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Invalid username or password")
		return
	}
	expirationTime := time.Now().Add(15 * time.Minute)
	tokenString, err := lib.CreateJWT(login.Username, expirationTime)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to create JWT")
		return
	}
	createApiResponse(w, lib.Token{
		Token:          tokenString,
		ExpirationTime: expirationTime.UTC().String(),
	}, http.StatusOK, "success", "none")
}
func (d *Handler) Register(w http.ResponseWriter, r *http.Request) {
	cred := models.Credentials{}
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to decode json")
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cred.Password), 8)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to create user account")
		return
	}
	user := &models.User{
		ID:       uuid.NewV4().String(),
		Email:    cred.Email,
		Username: cred.Username,
		Password: string(hashedPassword),
	}
	if err := models.RegisterUser(d.Client, d.DatabaseName, user); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to create user account")
		return
	}
	createApiResponse(w, nil, http.StatusOK, "success", "none")
}
func (d *Handler) Renew(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("token").(jwtBody)
	newToken, expTime, err := lib.RenewJWT(w, r, tkn.Token)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to renew JWT")
		return
	}
	createApiResponse(w, lib.Token{
		Token:          newToken,
		ExpirationTime: expTime,
	}, http.StatusOK,
		"success",
		"none")
}
func (d *Handler) GetUserinfo(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("token").(jwtBody)
	user, err := models.GetUserinfo(d.Client, d.DatabaseName, tkn.Username)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to return userinfo")
		return
	}
	createApiResponse(w, user, http.StatusOK, "success", "none")
}
