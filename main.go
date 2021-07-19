package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Daniorocket/cookit-backend/routing"
)

func main() {
	port := os.Getenv("PORT")
	router, err := routing.NewRouter()
	if err != nil {
		log.Print("Failed to init Router", err)
		return
	}
	fmt.Println("Server started!\nListening on :" + port)
	if err = http.ListenAndServe(":"+port, router); err != nil {
		log.Println("Failed to close server: ", err)
		return
	}
}
