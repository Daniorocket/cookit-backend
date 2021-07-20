package lib

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

type File struct {
	EncodedURL string `json:"encodedURL" bson:"encoded_url"`
	Extension  string `json:"extension" bson:"extension"`
}

func DecodeMultipartRequest(r *http.Request, data interface{}) (*multipart.Part, error) {

	var file *multipart.Part
	var js *json.Decoder
	mr, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF { //End of multipart data
			break
		}
		if err != nil {
			return nil, err
		}
		if part.FormName() == "file" {
			file = part
		}
		if part.FormName() == "json" {
			js = json.NewDecoder(part)
			if err := js.Decode(&data); err != nil {
				return nil, err
			}
		}
	}

	return file, nil
}
