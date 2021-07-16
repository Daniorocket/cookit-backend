package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"

	"github.com/Daniorocket/cookit-backend/models"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func (d *Handler) GetListOfCategories(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	categories, te, err := models.GetAllCategories(d.Client, d.DatabaseName, page, limit)
	if err != nil {
		log.Println("Error:", err)
		if err := CreateApiResponse(w, nil, http.StatusBadRequest, "failed", "Failed to get list of categories:"); err != nil {
			log.Println("Error create Api response:", err)
		}
		return
	}
	if err := CreateApiResponse(w, Pagination{
		Data:          categories,
		Limit:         limit,
		Page:          page,
		TotalElements: te,
	}, http.StatusOK, "success", "none"); err != nil {
		log.Println("Error create Api response:", err)
	}
}
func (d *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	mr, err := r.MultipartReader()
	if err != nil {
		log.Println("Error:", err)
		if err := CreateApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to load image:"); err != nil {
			log.Println("Error create Api response:", err)
		}
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
			if err := CreateApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to read json:"); err != nil {
				log.Println("Error create Api response:", err)
			}
			return
		}
		if part.FormName() == "file" {
			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, part); err != nil {
				log.Println("Error:", err)
				return
			}
			cat.File.EncodedURL = base64.StdEncoding.EncodeToString(buf.Bytes())
			cat.File.Extension = path.Ext(part.FileName())
			cat.ID = uuid.NewV4().String()
			switch ext := cat.File.Extension; ext {
			case ".jpg", ".JPG", ".png", ".PNG":
			default:
				log.Println("Error extension image")
				if err := CreateApiResponse(w, nil, http.StatusInternalServerError, "failed", "Bad image extension:"); err != nil {
					log.Println("Error create Api response:", err)
				}
				return
			}
		}
		if part.FormName() == "json" {
			jsonDecoder := json.NewDecoder(part)
			err = jsonDecoder.Decode(&cat)
			if err != nil {
				log.Println("Error:", err)
				if err := CreateApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to decode json:"); err != nil {
					log.Println("Error create Api response:", err)
				}
				return
			}
		}
	}
	if err := models.CreateCategory(d.Client, d.DatabaseName, cat); err != nil {
		log.Println("Error:", err)
		if err := CreateApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to create category:"); err != nil {
			log.Println("Error create Api response:", err)
		}
		return
	}
	if err := CreateApiResponse(w, nil, http.StatusOK, "success", "none"); err != nil {
		log.Println("Error create Api response:", err)
	}
}
func (d *Handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	category, err := models.GetCategoryByID(d.Client, d.DatabaseName, id)
	if err != nil {
		log.Println("Error:", err)
		if err := CreateApiResponse(w, nil, http.StatusInternalServerError, "failed", "Failed to get category:"); err != nil {
			log.Println("Error create Api response:", err)
		}
		return
	}
	if err := CreateApiResponse(w, category, http.StatusOK, "success", "none"); err != nil {
		log.Println("Error create Api response:", err)
	}
}
