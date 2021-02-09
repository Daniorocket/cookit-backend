package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Daniorocket/cookit-backend/models"
	uuid "github.com/satori/go.uuid"
)

func (d *Handler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateRecipe")
	recipe := &models.Recipe{}
	recipe.UserID = uuid.NewV4().String()
	recipe.ID = uuid.NewV4().String()
	fmt.Println("Recipe:", recipe)
	err := json.NewDecoder(r.Body).Decode(recipe)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Failed to decode body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := models.CreateRecipe(d.Sess, d.DatabaseName, *recipe); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
}
func (d *Handler) ListRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := models.GetAllRecipes(d.Sess, d.DatabaseName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	if err != nil {
		log.Println("Failed to prepare json describe list of users: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Pagination{
		Data:          recipes,
		Limit:         1,
		Page:          1,
		TotalElements: 3,
	})
}
