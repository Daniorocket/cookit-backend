package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Daniorocket/cookit-backend/db"
	"github.com/Daniorocket/cookit-backend/routing"
)

func main() {
	port := os.Getenv("PORT")
	categoryRepository, recipeRepository, authRepository, err := db.InitMongoDatabase()
	if err != nil {
		log.Print("Failed to init database:", err)
		return
	}
	router, err := routing.NewRouter(categoryRepository, recipeRepository, authRepository)
	if err != nil {
		log.Print("Failed to init router:", err)
		return
	}
	fmt.Println("Server started!\nListening on :" + port)
	if err = http.ListenAndServe(":"+port, router); err != nil {
		log.Println("Failed to close server: ", err)
		return
	}
}
