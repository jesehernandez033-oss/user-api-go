package controllers

import (
	"awesomeProject1/domain"
	"awesomeProject1/infrastructure"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User

	// 1. Leer JSON
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		infrastructure.Logger.Printf(
			"ERROR | endpoint=UpdateUser | error=invalid JSON | err=%v",
			err,
		)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Email obligatorio
	if user.Email == "" {
		infrastructure.Logger.Printf(
			"WARNING | endpoint=UpdateUser | missing email",
		)
		http.Error(w, "email is required to update user", http.StatusBadRequest)
		return
	}

	// 3. Verificar existencia
	var exists int
	err = infrastructure.DB.QueryRow(
		"SELECT COUNT(*) FROM users WHERE email = @p1",
		user.Email,
	).Scan(&exists)

	if err != nil {
		infrastructure.Logger.Printf(
			"ERROR | endpoint=UpdateUser | email=%s | db_error=%v",
			user.Email,
			err,
		)
		http.Error(w, "error checking user", http.StatusInternalServerError)
		return
	}

	if exists == 0 {
		infrastructure.Logger.Printf(
			"WARNING | endpoint=UpdateUser | user not found | email=%s",
			user.Email,
		)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// 4. Validaciones

	if user.Email != "" && !user.IsValidEmail() {
		infrastructure.Logger.Printf(
			"WARNING | endpoint=UpdateUser | invalid email format | email=%s",
			user.Email,
		)
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}

	if user.Password != "" {
		if !user.IsValidPassword() {
			infrastructure.Logger.Printf(
				"WARNING | endpoint=UpdateUser | weak password | email=%s",
				user.Email,
			)
			http.Error(w, "invalid password", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			infrastructure.Logger.Printf(
				"ERROR | endpoint=UpdateUser | email=%s | bcrypt_error=%v",
				user.Email,
				err,
			)
			http.Error(w, "error encrypting password", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)
	}

	if user.Birthday != "" {
		if !user.IsAdult() {
			infrastructure.Logger.Printf(
				"WARNING | endpoint=UpdateUser | underage update attempt | email=%s",
				user.Email,
			)
			http.Error(w, "user must be older than 14", http.StatusBadRequest)
			return
		}
	}

	// 🔥 5. QUERY DINÁMICA
	query := "UPDATE users SET "
	params := []interface{}{}
	i := 1

	updatedFields := []string{} // 👈 para logging

	if user.UserName != "" {
		query += "username = @p" + fmt.Sprint(i) + ", "
		params = append(params, user.UserName)
		updatedFields = append(updatedFields, "username")
		i++
	}

	if user.Password != "" {
		query += "password = @p" + fmt.Sprint(i) + ", "
		params = append(params, user.Password)
		updatedFields = append(updatedFields, "password")
		i++
	}

	if user.PhoneNumber != "" {
		query += "phone_number = @p" + fmt.Sprint(i) + ", "
		params = append(params, user.PhoneNumber)
		updatedFields = append(updatedFields, "phone_number")
		i++
	}

	if user.Birthday != "" {
		query += "birthday = @p" + fmt.Sprint(i) + ", "
		params = append(params, user.Birthday)
		updatedFields = append(updatedFields, "birthday")
		i++
	}

	if user.Address != "" {
		query += "address = @p" + fmt.Sprint(i) + ", "
		params = append(params, user.Address)
		updatedFields = append(updatedFields, "address")
		i++
	}

	// ⚠️ evitar update vacío
	if len(params) == 0 {
		infrastructure.Logger.Printf(
			"WARNING | endpoint=UpdateUser | no fields to update | email=%s",
			user.Email,
		)
		http.Error(w, "no fields to update", http.StatusBadRequest)
		return
	}

	query = strings.TrimSuffix(query, ", ")
	query += " WHERE email = @p" + fmt.Sprint(i)
	params = append(params, user.Email)

	// ejecutar
	_, err = infrastructure.DB.Exec(query, params...)
	if err != nil {
		infrastructure.Logger.Printf(
			"ERROR | endpoint=UpdateUser | email=%s | update_error=%v",
			user.Email,
			err,
		)
		http.Error(w, "error updating user", http.StatusInternalServerError)
		return
	}

	// 🟢 Log de éxito (sin datos sensibles)
	infrastructure.Logger.Printf(
		"INFO | user updated | email=%s | fields=%v",
		user.Email,
		updatedFields,
	)

	// 6. respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user updated successfully",
	})
}