package controllers

import (
	"awesomeProject1/application"
	"awesomeProject1/infrastructure"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponseToken struct {
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
}

// 🔐 Rate limit en memoria
var loginAttempts = make(map[string]int)
var lastAttempt = make(map[string]time.Time)
var mu sync.Mutex

func isBlocked(ip string) bool {
	mu.Lock()
	defer mu.Unlock()

	// Reset cada minuto
	if time.Since(lastAttempt[ip]) > time.Minute {
		loginAttempts[ip] = 0
	}

	if loginAttempts[ip] >= 5 {
		infrastructure.Logger.Printf(
			"WARNING | rate limit exceeded | ip=%s",
			ip,
		)
		return true
	}

	loginAttempts[ip]++
	lastAttempt[ip] = time.Now()

	return false
}

// Login godoc
// @Summary Iniciar sesión
// @Description Valida credenciales y devuelve un token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Credenciales"
// @Success 200 {object} ResponseToken
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 429 {string} string "Too many requests"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// 🔐 Obtener IP sin puerto
	ip := strings.Split(r.RemoteAddr, ":")[0]

	// 1. Leer JSON
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		infrastructure.Logger.Printf(
			"ERROR | endpoint=Login | ip=%s | error=invalid JSON | err=%v",
			ip,
			err,
		)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Validación básica
	if req.Email == "" || req.Password == "" {
		infrastructure.Logger.Printf(
			"WARNING | endpoint=Login | ip=%s | missing fields",
			ip,
		)
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	// 🔐 Bloqueo por intentos
	if isBlocked(ip) {
		http.Error(w, "too many attempts, try later", http.StatusTooManyRequests)
		return
	}

	// 3. Buscar usuario en DB
	var storedPassword string

	err = infrastructure.DB.QueryRow(
		"SELECT password FROM users WHERE email = @p1",
		req.Email,
	).Scan(&storedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			infrastructure.Logger.Printf(
				"WARNING | login failed | email=%s | ip=%s | reason=user not found",
				req.Email,
				ip,
			)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		infrastructure.Logger.Printf(
			"ERROR | endpoint=Login | email=%s | ip=%s | db_error=%v",
			req.Email,
			ip,
			err,
		)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// 4. Comparar password
	err = bcrypt.CompareHashAndPassword(
		[]byte(storedPassword),
		[]byte(req.Password),
	)

	if err != nil {
		infrastructure.Logger.Printf(
			"WARNING | login failed | email=%s | ip=%s | reason=wrong password",
			req.Email,
			ip,
		)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// ✅ Reset intentos si login exitoso
	mu.Lock()
	loginAttempts[ip] = 0
	mu.Unlock()

	// 5. Generar token
	token := application.GenerateTokenWithEmail(req.Email)

	// 🟢 Login exitoso (log informativo)
	infrastructure.Logger.Printf(
		"INFO | login successful | email=%s | ip=%s",
		req.Email,
		ip,
	)

	// 6. Respuesta
	response := ResponseToken{
		Message:     "login successful",
		AccessToken: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
