package controllers

import (
	"awesomeProject1/infrastructure"
	"encoding/json"
	"net/http"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	// 1. Obtener email desde query param
	email := r.URL.Query().Get("email")

	if email == "" {
		http.Error(w, "email is required", 400)
		return
	}

	// 2. Verificar si el usuario existe
	var exists int
	err := infrastructure.DB.QueryRow(
		"SELECT COUNT(*) FROM users WHERE email = @p1",
		email,
	).Scan(&exists)

	if err != nil {
		http.Error(w, "error checking user", 500)
		return
	}

	if exists == 0 {
		http.Error(w, "user not found", 404)
		return
	}

	// 3. Eliminar usuario
	_, err = infrastructure.DB.Exec(
		"DELETE FROM users WHERE email = @p1",
		email,
	)

	if err != nil {
		http.Error(w, "error deleting user", 500)
		return
	}

	// 4. Respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user deleted successfully",
	})
}
