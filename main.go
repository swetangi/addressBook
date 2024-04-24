package main

import (
	"addressBook/config"
	"addressBook/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	appCtx, err := config.NewAppCtx()
	if err != nil {
		log.Println("App ctx failed", err)
	}
	router := mux.NewRouter()
	router = routes.RouterList(appCtx, router)
	fmt.Println("Server Listening on Port No. 8080")
	// http.ListenAndServe("localhost:8080", router)
	http.ListenAndServe("localhost:8080",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
		)(router))

}
