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

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name UserProvider
type UserProvider interface {
	CreateUser(email, role, password string) error
}

// @Summary User Creation
// @Description Creates a new user with the provided email, role, and password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body models.UserCreate true "User creation details"
// @Success 200 {string} string "Successfully created a new user"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /create/user [post]
func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	const op = "handler.createUser"

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

	if user.Email == "" {
		h.log.Error("email is empty")
		http.Error(w, "email is empty", http.StatusBadRequest)
		return
	}
	if user.Role == "" {
		h.log.Error("role is empty")
		http.Error(w, "role is empty", http.StatusBadRequest)
		return
	}
	if user.Password == "" {
		h.log.Error("password is empty")
		http.Error(w, "password is empty", http.StatusBadRequest)
		return
	}

	err = h.userProvider.CreateUser(user.Email, user.Role, user.Password)
	if err != nil {
		log.Error("failed to create user", sl.Err(err))
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Created user with email - " + user.Email))
}
