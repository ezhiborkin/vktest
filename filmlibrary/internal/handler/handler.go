package handler

import (
	"encoding/json"
	"filmlibrary/internal/domain/models"
	"fmt"
	"net/http"
)

type Handler struct {
	serviceCooker ServiceCooker
}

type ServiceCooker interface {
	EditActor(actor *models.Actor) error
	AddActor(actor *models.Actor) error
	EditMovie(movie *models.Movie) error
	AddMovie(movie *models.Movie) error
}

func New(serviceCooker ServiceCooker) *Handler {
	return &Handler{serviceCooker: serviceCooker}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/edit/actor", onlyPostMiddleware(h.editActor))
	mux.HandleFunc("/edit/movie", onlyPostMiddleware(h.editMovie))
	mux.HandleFunc("/create/actor", onlyPostMiddleware(h.addActor))
	mux.HandleFunc("/create/movie", onlyPostMiddleware(h.addMovie))

	return mux
}

func (h *Handler) editActor(w http.ResponseWriter, r *http.Request) {
	const op = "handler.editActor"

	var req models.Actor
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, op, http.StatusBadRequest)
	}

	fmt.Println(&req.Name, &req.Sex)

	err = h.serviceCooker.EditActor(&req)
	if err != nil {
		http.Error(w, op, http.StatusNotImplemented)
	}

	response := map[string]interface{}{
		"kekys": 100,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) addActor(w http.ResponseWriter, r *http.Request) {
	const op = "handler.addActor"

	var req models.Actor
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, op, http.StatusBadRequest)
	}

	fmt.Println(req.Birthday)

	err = h.serviceCooker.AddActor(&req)
	if err != nil {
		http.Error(w, op, http.StatusNotImplemented)
	}

	response := map[string]interface{}{
		"kekys": 100,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) addMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.addMovie"

	var req models.Movie
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, op, http.StatusBadRequest)
	}

	err = h.serviceCooker.AddMovie(&req)
	if err != nil {
		http.Error(w, op, http.StatusNotImplemented)
	}

	response := map[string]interface{}{
		"kekys": 100,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) editMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.editActor"

	var req models.Movie
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, op, http.StatusBadRequest)
	}

	fmt.Println(req)

	err = h.serviceCooker.EditMovie(&req)
	if err != nil {
		http.Error(w, op, http.StatusNotImplemented)
	}

	response := map[string]interface{}{
		"movie": 100,
	}

	json.NewEncoder(w).Encode(response)
}

func onlyGetMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "only GET method is allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func onlyPostMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	}
}
