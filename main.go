package main

import (
	"awesomeProject1/application/controllers"
	"awesomeProject1/application/middleware"
	"awesomeProject1/infrastructure"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	routerParent := mux.NewRouter().StrictSlash(true)
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})
	infrastructure.ConnectDB()
	routerParent.HandleFunc("/ping", controllers.Ping).Methods("GET")
	routerParent.HandleFunc("/login", controllers.Login).Methods("POST")
	routerParent.HandleFunc("/uploadFile", controllers.UploadFile).Methods("POST")
	routerParent.HandleFunc("/createUser", controllers.CreateUser).Methods("POST")
	//routerParent.HandleFunc("/user", controllers.Login).Methods("GET")
	//routerParent.HandleFunc("/refresh", controllers.Login).Methods("POST")
	routerParent.HandleFunc("/updateUser", middleware.AuthMiddleware(controllers.UpdateUser)).Methods("PUT")
	routerParent.HandleFunc("/deleteUser", middleware.AuthMiddleware(controllers.DeleteUser)).Methods("DELETE")

	handler := corsMiddleware.Handler(routerParent)

	http.ListenAndServe(":8080", handler)
}
