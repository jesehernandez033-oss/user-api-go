package main

import (
	_ "awesomeProject1/docs"

	"awesomeProject1/application/controllers"
	"awesomeProject1/application/middleware"
	"awesomeProject1/infrastructure"

	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title User API Go
// @version 2.0
// @description API REST para gestión de usuarios con JWT, bcrypt, validaciones, logging y SQL Server.
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	routerParent := mux.NewRouter().StrictSlash(true)

	// 🔐 Obtener origen permitido
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
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

	// 🔥 INIT LOGGER
	infrastructure.InitLogger()

	// 🔥 LOG DE INICIO
	infrastructure.Logger.Println("Starting server...")

	// DB
	infrastructure.ConnectDB()

	// 📡 Routes
	routerParent.HandleFunc("/ping", controllers.Ping).Methods("GET")
	routerParent.HandleFunc("/login", controllers.Login).Methods("POST")
	routerParent.HandleFunc("/uploadFile", controllers.UploadFile).Methods("POST")
	routerParent.HandleFunc("/createUser", controllers.CreateUser).Methods("POST")
	routerParent.HandleFunc("/updateUser", middleware.AuthMiddleware(controllers.UpdateUser)).Methods("PUT")
	routerParent.HandleFunc("/deleteUser", middleware.AuthMiddleware(controllers.DeleteUser)).Methods("DELETE")

	// 📚 Swagger
	routerParent.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	handler := corsMiddleware.Handler(routerParent)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	infrastructure.Logger.Println("Server running on port:", port)

	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatal(err)
	}
}
