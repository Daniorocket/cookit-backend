package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/Daniorocket/cookit-backend/models"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func (d *Handler) GetListOfCategories(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	p, err := strconv.Atoi(page)
	if err != nil {
	}
	l, err := strconv.Atoi(limit)
	if err != nil {
	}

	categories, te, err := models.GetAllCategories(d.Client, d.DatabaseName, p, l)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to get list of categories:")
		return
	}

	createApiResponse(w, paginationResponse{
		Data:          categories,
		Limit:         limit,
		Page:          page,
		TotalElements: te,
	}, http.StatusOK, "success", "none")
}
func (d *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	mr, err := r.MultipartReader()
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to load image:")
		return
	}

	cat := models.Category{}
	for {
		part, err := mr.NextPart()
		if err == io.EOF { //End of multipart data
			break
		}
		if err != nil {
			log.Println("Error:", err)
			createApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to read json:")
			return
		}
		if part.FormName() == "file" {
			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, part); err != nil {
				log.Println("Error:", err)
				createApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to read file:")
				return
			}
			cat.File.EncodedURL = base64.StdEncoding.EncodeToString(buf.Bytes())
			cat.File.Extension = path.Ext(part.FileName())
			cat.ID = uuid.NewV4().String()
			switch ext := cat.File.Extension; ext {
			case ".jpg", ".JPG", ".png", ".PNG":
			default:
				log.Println("Error extension image")
				createApiResponse(w, nil, http.StatusInternalServerError, "failed", "Bad image extension:")
				return
			}
		}
		if part.FormName() == "json" {
			jsonDecoder := json.NewDecoder(part)
			err = jsonDecoder.Decode(&cat)
			if err != nil {
				log.Println("Error:", err)
				createApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to decode json:")
				return
			}
		}
	}
	if err := models.CreateCategory(d.Client, d.DatabaseName, cat); err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to create category:")
		return
	}
	createApiResponse(w, nil, http.StatusOK, "success", "none")
}
func (d *Handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	category, err := models.GetCategoryByID(d.Client, d.DatabaseName, id)
	if err != nil {
		log.Println("Error:", err)
		createApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to get category:")
		return
	}
	createApiResponse(w, category, http.StatusOK, "success", "none")
}
