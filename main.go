package main

import (
	"awesomeProject1/application/controllers"
	"awesomeProject1/application/middleware"
	"awesomeProject1/infrastructure"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	routerParent := mux.NewRouter().StrictSlash(true)

	// 🔐 Obtener origen permitido
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000" // fallback para desarrollo
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{allowedOrigin},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	// 🔥 INIT LOGGER (AQUÍ ESTÁ LA CLAVE)
	infrastructure.InitLogger()

	// 🔥 LOG DE INICIO
	infrastructure.Logger.Println("Starting server...")

	// DB
	infrastructure.ConnectDB()

	routerParent.HandleFunc("/ping", controllers.Ping).Methods("GET")
	routerParent.HandleFunc("/login", controllers.Login).Methods("POST")
	routerParent.HandleFunc("/uploadFile", controllers.UploadFile).Methods("POST")
	routerParent.HandleFunc("/createUser", controllers.CreateUser).Methods("POST")
	routerParent.HandleFunc("/updateUser", middleware.AuthMiddleware(controllers.UpdateUser)).Methods("PUT")
	routerParent.HandleFunc("/deleteUser", middleware.AuthMiddleware(controllers.DeleteUser)).Methods("DELETE")

	handler := corsMiddleware.Handler(routerParent)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	infrastructure.Logger.Println("Server running on port:", port)

	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatal(err) // este mata la app si falla al iniciar
	}
}