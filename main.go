package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Daniorocket/cookit-backend/handlers"
	"github.com/Daniorocket/cookit-backend/routing"
)

// let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database

func main() {
	router, err := routing.NewRouter()
	if err != nil {
		log.Print("Failed to init Router", err)
		return
	}
	router.Use(handlers.Authenticate)

	fmt.Println("Server started!")
	if err = http.ListenAndServe(":8080", router); err != nil {
		log.Println("Failed to close server: ", err)
		return
	}
}
