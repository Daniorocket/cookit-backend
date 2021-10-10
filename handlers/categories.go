package handlers

import (
	"log"
	"net/http"

	"github.com/Daniorocket/cookit-backend/lib"
	"github.com/Daniorocket/cookit-backend/models"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/validator.v2"
)

var errorGetCategory = "Nie można pobrać kategorii"
var errorCreateCategory = "Nie można utworzyć kategorii z wprowadzonymi danymi"

func (d *Handler) GetListOfCategories(w http.ResponseWriter, r *http.Request) {

	page, limit, err := lib.GetPageAndLimitFromRequest(r)
	if err != nil {
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", errorGetCategory)
		log.Println("Error:", err)
		return
	}
	categories, te, err := d.CategoryRepository.GetAll(page, limit)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorGetCategory)
		return
	}

	createApiResponse(w, paginationResponse{
		Data:          categories,
		Limit:         limit,
		Page:          page,
		TotalElements: te,
	}, http.StatusOK, "success", noError)
}
func (d *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {

	cat := models.Category{}
	encFile, ext, err := lib.DecodeMultipartRequest(r, &cat)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", errorMultipartData)
		return
	}

	if err := validator.Validate(cat); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", errorJSON)
		return
	}

	cat.File.EncodedURL = encFile
	cat.File.Extension = ext
	cat.ID = uuid.NewV4().String()

	if err := d.CategoryRepository.Create(cat); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", errorCreateCategory)
		return
	}
	createApiResponse(w, nil, http.StatusOK, "success", noError)
}
func (d *Handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	category, err := d.CategoryRepository.GetByID(id)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", errorGetCategory)
		return
	}
	createApiResponse(w, category, http.StatusOK, "success", noError)
}
