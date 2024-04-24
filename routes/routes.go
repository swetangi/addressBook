package routes

import (
	"addressBook/config"
	"addressBook/controllers/contacts"
	"addressBook/controllers/users"
	"addressBook/middleware"

	"github.com/gorilla/mux"
)

func RouterList(appctx *config.AppCtx, router *mux.Router) *mux.Router {

	userController := users.NewUserController(appctx)

	contactController := contacts.NewContactController(appctx)
	router.HandleFunc("/users/register", userController.RegisterUser).Methods("POST")
	router.HandleFunc("/users/login", userController.LoginUser).Methods("POST")

	route := router.PathPrefix("/user").Subrouter()
	route.Use(middleware.AuthMiddleware)
	route.HandleFunc("/contact", (contactController.CreateContact)).Methods("POST")
	route.HandleFunc("/contacts", contactController.GetContacts).Methods("GET")
	route.HandleFunc("/contacts/{cid}", (contactController.UpdateContact)).Methods("PATCH")
	route.HandleFunc("/contact/{cid}", (contactController.DeleteContact)).Methods("DELETE")

	route.HandleFunc("/contact/download", contactController.DownloadCSV).Methods("POST")
	route.HandleFunc("/logout", userController.LogoutUser).Methods("POST")
	return router
}
