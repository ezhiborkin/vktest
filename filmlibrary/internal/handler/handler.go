package handler

import (
	"encoding/json"
	"filmlibrary/internal/domain/models"
	"fmt"
	"net/http"
	"strconv"
)

type Handler struct {
	serviceCooker ServiceCooker
}

type ServiceCooker interface {
	EditActor(actor *models.Actor) error
	AddActor(actor *models.Actor) error
	GetActorsStorage() ([]*models.ActorListing, error)
	EditMovie(movie *models.Movie) error
	AddMovie(movie *models.Movie) error
	DeleteActor(id int64) error
	DeleteMovie(id int64) error
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

	mux.HandleFunc("/delete/actor", onlyDeleteMiddleware(h.deleteActor))
	mux.HandleFunc("/delete/movie", onlyDeleteMiddleware(h.deleteMovie))

	mux.HandleFunc("/get/actors", onlyGetMiddleware(h.getActors))

	return mux
}

func (h *Handler) editActor(w http.ResponseWriter, r *http.Request) {
	const op = "handler.editActor"

	var req models.Actor
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, op, http.StatusBadRequest)
	}

	err = h.serviceCooker.EditActor(&req)
	if err != nil {
		http.Error(w, op, http.StatusNotImplemented)
	}

	//TODO

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getActors(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getActors"

	actors, err := h.serviceCooker.GetActorsStorage()
	if err != nil {
		http.Error(w, "Failed to fetch actors", http.StatusInternalServerError)
		fmt.Printf("%s: %w", op, err)
		return
	}

	actorJSON, err := json.Marshal(actors)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		fmt.Printf("%s: %v\n", op, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(actorJSON)
}

func (h *Handler) deleteActor(w http.ResponseWriter, r *http.Request) {
	const op = "handler.deleteActor"

	actorIDStr := r.URL.Query().Get("id")
	actorID, err := strconv.ParseInt(actorIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}

	err = h.serviceCooker.DeleteActor(actorID)
	if err != nil {
		http.Error(w, op, http.StatusNotImplemented)
	}

	//TODO

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) addActor(w http.ResponseWriter, r *http.Request) {
	const op = "handler.addActor"

	var req models.Actor
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, op, http.StatusBadRequest)
	}

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

func (h *Handler) deleteMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.deleteMovie"

	movieIDStr := r.URL.Query().Get("id")
	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	err = h.serviceCooker.DeleteMovie(movieID)
	if err != nil {
		//log
		http.Error(w, op, http.StatusNotImplemented)
	}

	//TODO

	w.WriteHeader(http.StatusOK)
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

func onlyDeleteMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	}
}
