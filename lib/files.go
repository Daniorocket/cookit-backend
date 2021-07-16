package lib

type File struct {
	EncodedURL string `json:"encodedURL" bson:"encoded_url"`
	Extension  string `json:"extension" bson:"extension"`
}
