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

var errorUsername = "Wprowadzono nieprawidłowy login lub hasło."
var errorCreateJWT = "Nie można utworzyć tokenu uwierzytelniającego."
var errorCreateAccount = "Nazwa użytkownika jest nieprawidłowa, lub już została zajerestrowana w systemie."
var errorVerifyEmail = "Email jest nieprawidłowy."
var errorCheckEmail = "Jeśli wprowadzony adres e-mail jest prawidłowy, to proszę sprawdzić skrzynkę odbiorczą."
var errorResetPassword = "Wprowadzone hasło jest nieprawidłowe."
var errorSendEmail = "Nie udało się wysłać e-maila."

func (d *Handler) Login(w http.ResponseWriter, r *http.Request) {

	login := models.Login{}
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorJSON)
		return
	}

	passDB, err := d.AuthRepository.GetPassword(login.Username)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorUsername)
		return
	}

	if err := lib.CompareHashAndPassword([]byte(passDB), []byte(login.Password)); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorUsername)
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	tokenString, err := lib.CreateJWT(login.Username, expirationTime)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorCreateJWT)
		return
	}

	createApiResponse(w, lib.Token{
		Token:          tokenString,
		ExpirationTime: expirationTime.UTC().String(),
	},
		http.StatusOK,
		"success",
		noError)
}
func (d *Handler) Register(w http.ResponseWriter, r *http.Request) {

	cred := models.Credentials{}
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorJSON)
		return
	}

	hashedPassword, err := lib.HashPassword([]byte(cred.Password))
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorCreateAccount)
		return
	}

	user := models.User{
		ID:               uuid.NewV4().String(),
		Email:            cred.Email,
		Username:         cred.Username,
		Password:         string(hashedPassword),
		FavoritesRecipes: []models.Recipe{},
	}

	if err := d.AuthRepository.Register(user); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorCreateAccount)
		return
	}

	createApiResponse(w, nil, http.StatusOK, "success", noError)
}
func (d *Handler) Renew(w http.ResponseWriter, r *http.Request) {

	tkn, ok := r.Context().Value("token").(jwtBody)
	if !ok {
		log.Println("Error:", errors.New("Error parsing JWT"))
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorVerifyJWT)
		return
	}
	newToken, expTime, err := lib.RenewJWT(w, r, tkn.Token)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorVerifyJWT)
		return
	}

	createApiResponse(w, lib.Token{
		Token:          newToken,
		ExpirationTime: expTime,
	}, http.StatusOK,
		"success",
		noError)
}
func (d *Handler) GetUserinfo(w http.ResponseWriter, r *http.Request) {

	tkn, ok := r.Context().Value("token").(jwtBody)
	if !ok {
		log.Println("Error:", errors.New("Failed to read JWT from http header"))
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorVerifyJWT)
		return
	}

	user, err := d.AuthRepository.GetUserinfo(tkn.Username)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorUsername)
		return
	}

	createApiResponse(w, user, http.StatusOK, "success", noError)
}
func (d *Handler) RemindPassword(w http.ResponseWriter, r *http.Request) {
	email := models.Email{}
	if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorJSON)
		return
	}

	if err := validator.Validate(email); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorVerifyEmail)
		return
	}
	user, err := d.AuthRepository.CheckEmail(email.Email)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusOK, "success", errorCheckEmail)
		return
	}
	user.PasswordRemindID = uuid.NewV4().String()

	if err := d.AuthRepository.Update(user.ID, user); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorResetPassword)
		return
	}

	if err := lib.CreateEmail(email.Email, "Zmiana hasła w serwisie CookIT", "Dzień dobry.\nAby zmienić swoje hasło proszę przejść na stronę:"+
		"https://cookit0.herokuapp.com/api/v1/przypomnij-haslo/"+
		user.PasswordRemindID); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorSendEmail)
		return
	}
	createApiResponse(w, nil, http.StatusOK, "success", errorCheckEmail)
}
func (d *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	password := models.Password{}
	if err := json.NewDecoder(r.Body).Decode(&password); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorJSON)
		return
	}
	passwordRemindID := mux.Vars(r)["id"]
	user, err := d.AuthRepository.GetUserByPasswordRemindID(passwordRemindID)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorResetPassword)
		return
	}
	hashedPassword, err := lib.HashPassword([]byte(password.Password))
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "success", errorResetPassword)
		return
	}

	user.Password = string(hashedPassword)
	user.PasswordRemindID = ""
	if err := d.AuthRepository.Update(user.ID, user); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorResetPassword)
		return
	}
	createApiResponse(w, nil, http.StatusOK, "success", noError)
}
func (d *Handler) EditUserAccount(w http.ResponseWriter, r *http.Request) {
	tkn, ok := r.Context().Value("token").(jwtBody)
	if !ok {
		log.Println("Error:", errors.New("Error parsing JWT"))
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorVerifyJWT)
		return
	}
	usr, err := d.AuthRepository.GetUserinfo(tkn.Username)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorUsername)
		return
	}
	encFile, _, err := lib.DecodeMultipartRequest(r, &usr)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", errorMultipartData)
		return
	}
	if encFile != "" { //File is uploaded
		usr.AvatarURL = encFile
	}
	if err := d.AuthRepository.Update(usr.ID, usr); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorUpdateData)
		return
	}
	createApiResponse(w, nil, http.StatusOK, "success", noError)
}
