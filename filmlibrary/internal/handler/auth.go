package handler

import (
	"encoding/json"
	"errors"
	"filmlibrary/internal/domain/models"
	"filmlibrary/internal/lib/logger/sl"
	"io"
	"log/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name AuthProvider
type AuthProvider interface {
	LoginUser(email, password string) (string, error)
}

// @Summary User Login
// @Description Authenticates a user and generates an authentication token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body models.UserLogin true "User credentials for login"
// @Success 200 {string} string "Successfully logged in. Authentication token is included in the 'Authorization' header"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /login [post]
func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	const op = "handler.loginUser"

	log := h.log.With(slog.String("op", op))

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty", sl.Err(err))
		http.Error(w, "empty request", http.StatusBadRequest)
		return
	}

	log.Info("request body decoded")

	tokenString, err := h.authProvider.LoginUser(user.Email, user.Password)
	if err != nil {
		log.Error("failed to login user", sl.Err(err))
		http.Error(w, "failed to login user", http.StatusBadRequest)
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully logged in, your token: " + tokenString))
}
