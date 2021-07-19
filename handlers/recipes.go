package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Daniorocket/cookit-backend/models"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/validator.v2"
)

func (d *Handler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("token").(jwtBody)
	recipe := models.Recipe{
		ID:     uuid.NewV4().String(),
		UserID: tkn.Username,
		Date:   time.Now().UTC().String(),
	}
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Error:", err)
		if err := CreateApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to decode json"); err != nil {
			log.Println("Error create Api response:", err)
		}
		return
	}
	if err := validator.Validate(recipe); err != nil {
		log.Println("Error:", err)
		if err := CreateApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to validate json"); err != nil {
			log.Println("Error create Api response:", err)
		}
		return
	}
	if err := models.CreateRecipe(d.Client, d.DatabaseName, &recipe); err != nil {
		log.Println("Error:", err)
		if err := CreateApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to create recipe"); err != nil {
			log.Println("Error create Api response:", err)
		}
		return
	}
	CreateApiResponse(w, nil, http.StatusOK, "success", "none")
}
func (d *Handler) GetListOfRecipes(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	recipes, te, err := models.GetAllRecipes(d.Client, d.DatabaseName, page, limit)
	if err != nil {
		if err := CreateApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to get list of recipes"); err != nil {
			log.Println("Error create Api response:", err)
		}
		log.Println("Error:", err)
		return
	}
	if err := CreateApiResponse(w, paginationResponse{
		Data:          recipes,
		Limit:         limit,
		Page:          page,
		TotalElements: te,
	}, http.StatusOK, "success", "none"); err != nil {
		log.Println("Error create Api response:", err)
	}
}
func (d *Handler) GetListOfRecipesByTags(w http.ResponseWriter, r *http.Request) {
	data := make(map[string][]int)
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		if err := CreateApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to get list of recipes"); err != nil {
			log.Println("Error create Api response:", err)
		}
		log.Println("Error:", err)
		return
	}
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	recipes, te, err := models.GetAllRecipesByTags(d.Client, d.DatabaseName, data["tags"], page, limit)
	if err != nil {
		if err := CreateApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to get list of recipes from DB"); err != nil {
			log.Println("Error create Api response:", err)
		}
		log.Println("Error:", err)
		return
	}
	if err := CreateApiResponse(w, paginationResponse{
		Data:          recipes,
		Limit:         limit,
		Page:          page,
		TotalElements: te,
	}, http.StatusOK, "success", "none"); err != nil {
		log.Println("Error create Api response:", err)
	}
}
