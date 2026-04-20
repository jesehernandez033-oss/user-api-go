package controllers

import (
	"awesomeProject1/infrastructure"
	"encoding/json"
	"net/http"
	"time"
)

// Usa el mismo struct UserCredentials y la función isValidPassword
// que ya definiste en create_user.go

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user UserCredentials

	// 1. Leer JSON
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Email obligatorio (identificador)
	if user.Email == "" {
		http.Error(w, "email is required to update user", 400)
		return
	}

	// 3. Verificar que el usuario existe
	var exists int
	err = infrastructure.DB.QueryRow(
		"SELECT COUNT(*) FROM users WHERE email = @p1",
		user.Email,
	).Scan(&exists)

	if err != nil {
		http.Error(w, "error checking user", 500)
		return
	}

	if exists == 0 {
		http.Error(w, "user not found", 404)
		return
	}

	// 4. Validaciones (solo si vienen datos)

	// Password
	if user.Password != "" {
		if !isValidPassword(user.Password) {
			http.Error(w, "password must be 8-16 chars, include uppercase, number and special character", 400)
			return
		}
	}

	// Edad
	if user.Birthday != "" {
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
	}

	// 5. Actualizar en DB
	_, err = infrastructure.DB.Exec(
		`UPDATE users SET 
		 username = @p1,
		 password = @p2,
		 phone_number = @p3,
		 birthday = @p4,
		 address = @p5
		 WHERE email = @p6`,
		user.UserName,
		user.Password,
		user.PhoneNumber,
		user.Birthday,
		user.Address,
		user.Email,
	)

	if err != nil {
		http.Error(w, "error updating user", 500)
		return
	}

	// 6. Respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user updated successfully",
	})
}
