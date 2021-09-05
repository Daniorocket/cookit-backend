package lib

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"path"
)

type File struct {
	EncodedURL string `json:"encodedURL" bson:"encoded_url"`
	Extension  string `json:"extension" bson:"extension"`
}

func DecodeMultipartRequest(r *http.Request, data interface{}) (string, string, error) {
	buf := bytes.NewBuffer(nil)
	var js *json.Decoder
	var ext string
	mr, err := r.MultipartReader()
	if err != nil {
		return "", "", err
	}
	for {
		part, err := mr.NextPart()
		if err == io.EOF { //End of multipart data
			break
		}
		if err != nil {
			return "", "", err
		}
		if part.FormName() == "file" {
			if _, err := io.Copy(buf, part); err != nil {
				return "", "", err
			}
			ext = path.Ext(part.FileName())
			switch ext {
			case ".jpg", ".JPG", ".png", ".PNG":
			default:
				return "", "", err
			}
		}
		if part.FormName() == "json" {
			js = json.NewDecoder(part)
			if err := js.Decode(&data); err != nil {
				return "", "", err
			}
		}
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), ext, nil
}
