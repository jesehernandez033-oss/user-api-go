package controllers

import (
	"awesomeProject1/infrastructure"
	"encoding/json"
	"net/http"
)

type DeleteRequest struct {
	Email string `json:"email"`
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var req DeleteRequest

	// 1. Leer JSON
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		infrastructure.Logger.Printf(
			"ERROR | endpoint=DeleteUser | error=invalid JSON | err=%v",
			err,
		)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		infrastructure.Logger.Printf(
			"WARNING | endpoint=DeleteUser | missing email",
		)
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	// 2. Verificar existencia
	var exists int
	err = infrastructure.DB.QueryRow(
		"SELECT COUNT(*) FROM users WHERE email = @p1",
		req.Email,
	).Scan(&exists)

	if err != nil {
		infrastructure.Logger.Printf(
			"ERROR | endpoint=DeleteUser | email=%s | db_error=%v",
			req.Email,
			err,
		)
		http.Error(w, "error checking user", http.StatusInternalServerError)
		return
	}

	if exists == 0 {
		infrastructure.Logger.Printf(
			"WARNING | endpoint=DeleteUser | user not found | email=%s",
			req.Email,
		)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// 3. Eliminar
	_, err = infrastructure.DB.Exec(
		"DELETE FROM users WHERE email = @p1",
		req.Email,
	)

	if err != nil {
		infrastructure.Logger.Printf(
			"ERROR | endpoint=DeleteUser | email=%s | delete_error=%v",
			req.Email,
			err,
		)
		http.Error(w, "error deleting user", http.StatusInternalServerError)
		return
	}

	// 🟢 Log de éxito
	infrastructure.Logger.Printf(
		"INFO | user deleted | email=%s",
		req.Email,
	)

	// 4. Respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user deleted successfully",
	})
}