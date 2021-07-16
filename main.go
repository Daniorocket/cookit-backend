package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Daniorocket/cookit-back/handlers"
	"github.com/Daniorocket/cookit-back/routing"
)

func main() {
	port := os.Getenv("PORT")
	router, err := routing.NewRouter()
	if err != nil {
		log.Print("Failed to init Router", err)
		return
	}
	router.Use(handlers.Authenticate)
	fmt.Println("Server started!\nListening on :" + port)
	if err = http.ListenAndServe(":"+port, router); err != nil {
		log.Println("Failed to close server: ", err)
		return
	}
}
