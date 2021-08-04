package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Daniorocket/cookit-backend/lib"
	"github.com/Daniorocket/cookit-backend/models"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/validator.v2"
)

func (d *Handler) Login(w http.ResponseWriter, r *http.Request) {

	login := models.Login{}
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to decode json")
		return
	}

	passDB, err := d.AuthRepository.GetPassword(login.Username)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Invalid username or password")
		return
	}

	if err := lib.CompareHashAndPassword([]byte(passDB), []byte(login.Password)); err != nil {
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
	},
		http.StatusOK,
		"success",
		"none")
}
func (d *Handler) Register(w http.ResponseWriter, r *http.Request) {

	cred := models.Credentials{}
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to decode json")
		return
	}

	hashedPassword, err := lib.HashPassword([]byte(cred.Password))
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to create user account")
		return
	}

	user := models.User{
		ID:       uuid.NewV4().String(),
		Email:    cred.Email,
		Username: cred.Username,
		Password: string(hashedPassword),
	}

	if err := d.AuthRepository.Register(user); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to create user account")
		return
	}

	createApiResponse(w, nil, http.StatusOK, "success", "none")
}
func (d *Handler) Renew(w http.ResponseWriter, r *http.Request) {

	tkn, ok := r.Context().Value("token").(jwtBody)
	if !ok {
		log.Println("Error:", errors.New("Error parsing JWT"))
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to read JWT from http header")
		return
	}
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

	tkn, ok := r.Context().Value("token").(jwtBody)
	if !ok {
		log.Println("Error:", errors.New("Failed to read JWT from http header"))
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to read JWT from http header")
		return
	}

	user, err := d.AuthRepository.GetUserinfo(tkn.Username)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to return userinfo")
		return
	}

	createApiResponse(w, user, http.StatusOK, "success", "none")
}
func (d *Handler) RemindPassword(w http.ResponseWriter, r *http.Request) {
	email := models.Email{}
	if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to decode json")
		return
	}

	if err := validator.Validate(email); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to validate email")
		return
	}
	user, err := d.AuthRepository.CheckEmail(email.Email)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to remind password")
		return
	}
	user.PasswordRemindID = uuid.NewV4().String()

	if err := d.AuthRepository.Update(user); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to reset password")
		return
	}

	if err := lib.CreateEmail(email.Email, "Zmiana hasła w serwisie CookIT", "Dzień dobry.\nAby zmienić swoje hasło proszę przejść na stronę:"+
		"http://localhost:5000/api/v1/remindpassword/"+
		user.PasswordRemindID); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to send email")
		return
	}
	createApiResponse(w, nil, http.StatusOK, "success", "none")
}
func (d *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	password := models.Password{}
	if err := json.NewDecoder(r.Body).Decode(&password); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to decode json")
		return
	}
	passwordRemindID := mux.Vars(r)["id"]
	user, err := d.AuthRepository.GetUserByPasswordRemindID(passwordRemindID)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to reset password")
		return
	}
	hashedPassword, err := lib.HashPassword([]byte(password.Password))
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "success", "Failed to reset password")
		return
	}

	user.Password = string(hashedPassword)
	user.PasswordRemindID = ""
	if err := d.AuthRepository.Update(user); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to reset password")
		return
	}
	createApiResponse(w, nil, http.StatusOK, "failed", "none")
}
