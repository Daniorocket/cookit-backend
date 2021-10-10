package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Daniorocket/cookit-backend/lib"
	"github.com/Daniorocket/cookit-backend/models"
	"github.com/gorilla/mux"
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
	name := r.URL.Query().Get("name")
	categories := r.URL.Query().Get("categories")
	cat := strings.Split(categories, ",")
	page, limit, err := lib.GetPageAndLimitFromRequest(r)
	if err != nil {
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", errorGetRecipes)
		log.Println("Error:", err)
		return
	}
	recipes, te, err := d.RecipeRepository.GetAll(cat, page, limit, name)
	if err != nil {
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
func (d *Handler) GetRecipeByID(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	recipe, err := d.RecipeRepository.GetByID(id)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", errorGetCategory)
		return
	}
	createApiResponse(w, recipe, http.StatusOK, "success", noError)
}
func (d *Handler) AddToFavorites(w http.ResponseWriter, r *http.Request) {

	tkn := r.Context().Value("token").(jwtBody)
	id := mux.Vars(r)["id"]

	user, err := d.AuthRepository.GetUserinfo(tkn.Username)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorUsername)
		return
	}

	rec, err := d.RecipeRepository.GetByID(id)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", errorGetCategory)
		return
	}
	for _, v := range user.FavoritesRecipes {
		if v.ID == rec.ID {
			log.Println("Error:", errorRecipeExists)
			createApiResponse(w, nil, http.StatusBadRequest, "failed", errorUpdateData)
			return
		}
	}
	user.FavoritesRecipes = append(user.FavoritesRecipes, rec)
	if err := d.AuthRepository.Update(user.ID, user); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorUpdateData)
		return
	}

	createApiResponse(w, nil, http.StatusOK, "success", noError)
}
