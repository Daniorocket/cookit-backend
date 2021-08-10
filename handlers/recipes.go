package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Daniorocket/cookit-backend/lib"
	"github.com/Daniorocket/cookit-backend/models"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/validator.v2"
)

func (d *Handler) CreateRecipe(w http.ResponseWriter, r *http.Request) {

	tkn := r.Context().Value("token").(jwtBody)
	recipe := models.Recipe{
		ID:       uuid.NewV4().String(),
		Username: tkn.Username,
		Date:     time.Now().UTC().String(),
	}

	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to decode json")
		return
	}

	if err := validator.Validate(recipe); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to validate json")
		return
	}

	if err := d.RecipeRepository.Create(recipe); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to create recipe")
		return
	}
	createApiResponse(w, nil, http.StatusOK, "success", "none")
}
func (d *Handler) GetListOfRecipes(w http.ResponseWriter, r *http.Request) {

	page, limit, err := lib.GetPageAndLimitFromRequest(r)
	if err != nil {
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to get list of recipes")
		log.Println("Error:", err)
		return
	}
	recipes, te, err := d.RecipeRepository.GetAll(page, limit)
	if err != nil {
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to get list of recipes")
		log.Println("Error:", err)
		return
	}

	createApiResponse(w, paginationResponse{
		Data:          recipes,
		Limit:         limit,
		Page:          page,
		TotalElements: te,
	},
		http.StatusOK,
		"success",
		"none")
}
func (d *Handler) GetListOfRecipesByTags(w http.ResponseWriter, r *http.Request) {

	page, limit, err := lib.GetPageAndLimitFromRequest(r)
	if err != nil {
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to get list of recipes")
		log.Println("Error:", err)
		return
	}

	var tags models.TagsList
	if err := json.NewDecoder(r.Body).Decode(&tags); err != nil {
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to get list of recipes")
		log.Println("Error:", err)
		return
	}

	recipes, te, err := d.RecipeRepository.GetAllByTags(tags.Tags, page, limit)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to get list of recipes from DB")
		return
	}

	createApiResponse(w, paginationResponse{
		Data:          recipes,
		Limit:         limit,
		Page:          page,
		TotalElements: te,
	},
		http.StatusOK,
		"success",
		"none")
}
