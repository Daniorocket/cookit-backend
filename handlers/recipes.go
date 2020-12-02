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
	if err := models.CreateRecipe(d.Sess, d.DatabaseName, models.Recipe{
		ID:          uuid.NewV4().String(),
		Description: "Siema",
		Kitchen:     1,
		ListOfSteps: []string{},
		Name:        "Kuchnia polaka",
		Tags:        1,
		UserID:      uuid.NewV4().String(),
	}); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
}
func (d *Handler) ListRecipes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ListRecipes")

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
