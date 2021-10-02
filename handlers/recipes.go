package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Daniorocket/cookit-backend/lib"
	"github.com/Daniorocket/cookit-backend/models"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/validator.v2"
)

var errorCreateRecipe = "Nie można utworzyć przepisu z wprowadzonymi danymi."
var errorGetRecipes = "Nie można pobrać listy z przepisami."
var errorGetUnits = "Nie można pobrać listy z jednostkami."

func (d *Handler) CreateRecipe(w http.ResponseWriter, r *http.Request) {

	tkn := r.Context().Value("token").(jwtBody)

	recipe := models.Recipe{
		ID:       uuid.NewV4().String(),
		Username: tkn.Username,
		Date:     time.Now().UTC().String(),
	}
	encFile, ext, err := lib.DecodeMultipartRequest(r, &recipe)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", errorMultipartData)
		return
	}

	recipe.File.EncodedURL = encFile
	recipe.File.Extension = ext
	//Check if valid ingredients unitID
	for i, v := range recipe.Ingredients {
		if _, err := d.RecipeRepository.GetUnit(v.UnitID); err != nil {
			log.Println("Error:", err)
			createApiResponse(w, nil, http.StatusBadRequest, "failed", errorCreateRecipe)
			return
		}
		recipe.Ingredients[i].ID = uuid.NewV4().String()
	}

	//Check if valid categoriesID
	for _, v := range recipe.CategoriesID {
		if _, err := d.CategoryRepository.GetByID(v); err != nil {
			log.Println("Error:", err)
			createApiResponse(w, nil, http.StatusBadRequest, "failed", errorCreateRecipe)
			return
		}
	}

	if err := validator.Validate(recipe); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorJSON)
		return
	}

	if err := d.RecipeRepository.Create(recipe); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorCreateRecipe)
		return
	}
	createApiResponse(w, nil, http.StatusOK, "success", noError)
}
func (d *Handler) GetListOfRecipes(w http.ResponseWriter, r *http.Request) {

	page, limit, err := lib.GetPageAndLimitFromRequest(r)
	if err != nil {
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", errorGetRecipes)
		log.Println("Error:", err)
		return
	}
	recipes, te, err := d.RecipeRepository.GetAll(page, limit)
	if err != nil {
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", errorGetRecipes)
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
		noError)
}
func (d *Handler) GetListOfRecipesByCategories(w http.ResponseWriter, r *http.Request) {

	page, limit, err := lib.GetPageAndLimitFromRequest(r)
	if err != nil {
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", errorGetRecipes)
		log.Println("Error:", err)
		return
	}

	var categoriesID models.CategoryID
	if err := json.NewDecoder(r.Body).Decode(&categoriesID); err != nil {
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorGetRecipes)
		log.Println("Error:", err)
		return
	}

	fmt.Println("categories id:,", categoriesID)
	recipes, te, err := d.RecipeRepository.GetAllByCategories(categoriesID, page, limit)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorGetRecipes)
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
		noError)
}
func (d *Handler) GetUnits(w http.ResponseWriter, r *http.Request) {
	units, err := d.RecipeRepository.GetAllUnits()
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", errorGetUnits)
		return
	}
	createApiResponse(w, units, http.StatusOK, "success", noError)
}
