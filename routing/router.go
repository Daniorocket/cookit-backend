package routing

import (
	"log"
	"net/http"

	"github.com/Daniorocket/cookit-backend/handlers"
	"github.com/Daniorocket/cookit-backend/logger"
	"github.com/Daniorocket/cookit-backend/models"
	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
)

func NewRouter() (*mux.Router, error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal("Error connect to mongoDB")
		return nil, err
	}
	handler := handlers.Handler{
		Sess:         session,
		DatabaseName: "CookIt",
	}
	//Index for users
	keysUser := []string{"id", "email", "username"}
	for i := range keysUser {
		index := mgo.Index{
			Key:        []string{keysUser[i]},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     true,
		}
		if err := session.DB("CookIt").C(models.CollectionUsers).EnsureIndex(index); err != nil {
			return nil, err
		}
	}
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range initRoutes(handler) {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router, nil
}
