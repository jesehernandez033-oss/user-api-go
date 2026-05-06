package controllers

import (
	"awesomeProject1/domain"
	"awesomeProject1/infrastructure"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser godoc
// @Summary Crear usuario
// @Description Crea un nuevo usuario con validaciones y contraseña encriptada
// @Tags Users
// @Accept json
// @Produce json
// @Param user body domain.User true "Datos del usuario"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /createUser [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User

	// 1. Leer JSON
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		infrastructure.Logger.Printf(
			"ERROR | endpoint=CreateUser | error=invalid JSON | err=%v",
			err,
		)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Validaciones básicas
	if user.UserName == "" || user.Email == "" || user.Password == "" {
		infrastructure.Logger.Printf(
			"WARNING | endpoint=CreateUser | missing required fields | email=%s",
			user.Email,
		)
		http.Error(w, "username, email and password are required", http.StatusBadRequest)
		return
	}

	// 📧 Validaciones DOMAIN
	if !user.IsValidEmail() {
		infrastructure.Logger.Printf(
			"WARNING | endpoint=CreateUser | invalid email format | email=%s",
			user.Email,
		)
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}

	if !user.IsValidPassword() {
		infrastructure.Logger.Printf(
			"WARNING | endpoint=CreateUser | weak password | email=%s",
			user.Email,
		)
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	if !user.IsAdult() {
		infrastructure.Logger.Printf(
			"WARNING | endpoint=CreateUser | underage user | email=%s",
			user.Email,
		)
		http.Error(w, "user must be older than 14", http.StatusBadRequest)
		return
	}

	// 📧 Validar email único
	var exists int
	err = infrastructure.DB.QueryRow(
		"SELECT COUNT(*) FROM users WHERE email = @p1",
		user.Email,
	).Scan(&exists)

	if err != nil {
		infrastructure.Logger.Printf(
			"ERROR | endpoint=CreateUser | email=%s | db_error=%v",
			user.Email,
			err,
		)
		http.Error(w, "error checking email", http.StatusInternalServerError)
		return
	}

	if exists > 0 {
		infrastructure.Logger.Printf(
			"WARNING | endpoint=CreateUser | email already exists | email=%s",
			user.Email,
		)
		http.Error(w, "email already exists", http.StatusBadRequest)
		return
	}

	// 🔐 Encriptar password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		infrastructure.Logger.Printf(
			"ERROR | endpoint=CreateUser | email=%s | bcrypt_error=%v",
			user.Email,
			err,
		)
		http.Error(w, "error encrypting password", http.StatusInternalServerError)
		return
	}

	// 4. Insertar en DB
	_, err = infrastructure.DB.Exec(
		`INSERT INTO users 
		(username, email, password, phone_number, birthday, address) 
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6)`,
		user.UserName,
		user.Email,
		string(hashedPassword),
		user.PhoneNumber,
		user.Birthday,
		user.Address,
	)

	if err != nil {
		infrastructure.Logger.Printf(
			"ERROR | endpoint=CreateUser | email=%s | insert_error=%v",
			user.Email,
			err,
		)
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}

	// 🟢 Log de éxito
	infrastructure.Logger.Printf(
		"INFO | user created | email=%s",
		user.Email,
	)

	// 5. Respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user created successfully",
	})
}
