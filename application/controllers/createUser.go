package controllers

import (
	"awesomeProject1/infrastructure"
	"encoding/json"
	"net/http"
	"time"
	"unicode"
)

type UserCredentials struct {
	UserName    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Birthday    string `json:"birthday"`
	Address     string `json:"address"`
}

// 🔐 Validación de password SIN regex complejo
func isValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 16 {
		return false
	}

	var hasUpper, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasNumber = true
		case !unicode.IsLetter(char) && !unicode.IsDigit(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasNumber && hasSpecial
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user UserCredentials

	// 1. Leer JSON
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Validaciones

	// Campos obligatorios
	if user.UserName == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "username, email and password are required", 400)
		return
	}

	// Password válida
	if !isValidPassword(user.Password) {
		http.Error(w, "password must be 8-16 chars, include uppercase, number and special character", 400)
		return
	}

	// Validar fecha y edad
	birthDate, err := time.Parse("2006-01-02", user.Birthday)
	if err != nil {
		http.Error(w, "invalid date format (use YYYY-MM-DD)", 400)
		return
	}

	age := time.Now().Year() - birthDate.Year()
	if age < 14 {
		http.Error(w, "user must be older than 14", 400)
		return
	}

	// Validar email único
	var exists int
	err = infrastructure.DB.QueryRow(
		"SELECT COUNT(*) FROM users WHERE email = @p1",
		user.Email,
	).Scan(&exists)

	if err != nil {
		http.Error(w, "error checking email", 500)
		return
	}

	if exists > 0 {
		http.Error(w, "email already exists", 400)
		return
	}

	// 3. Insertar en DB
	_, err = infrastructure.DB.Exec(
		`INSERT INTO users 
		(username, email, password, phone_number, birthday, address) 
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6)`,
		user.UserName,
		user.Email,
		user.Password,
		user.PhoneNumber,
		user.Birthday,
		user.Address,
	)

	if err != nil {
		http.Error(w, "error creating user", 500)
		return
	}

	// 4. Respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user created successfully",
	})
}
