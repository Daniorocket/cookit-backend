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

	//Check if valid ingredients unitID
	for _, v := range recipe.Ingredients {
		if _, err := d.RecipeRepository.GetUnit(v.UnitID); err != nil {
			log.Println("Error:", err)
			createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to create recipe")
			return
		}
	}

	//Check if valid categoriesID
	for _, v := range recipe.CategoriesID {
		if _, err := d.CategoryRepository.GetByID(v); err != nil {
			log.Println("Error:", err)
			createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to create recipe")
			return
		}
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
func (d *Handler) GetListOfRecipesByCategories(w http.ResponseWriter, r *http.Request) {

	page, limit, err := lib.GetPageAndLimitFromRequest(r)
	if err != nil {
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to get list of recipes")
		log.Println("Error:", err)
		return
	}

	var categoriesID []string
	if err := json.NewDecoder(r.Body).Decode(&categoriesID); err != nil {
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to get list of categories")
		log.Println("Error:", err)
		return
	}

	recipes, te, err := d.RecipeRepository.GetAllByCategories(categoriesID, page, limit)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to get list of recipes")
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
