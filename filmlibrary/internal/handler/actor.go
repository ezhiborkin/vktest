package handler

import (
	"encoding/json"
	"filmlibrary/internal/domain/models"
	"filmlibrary/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"strconv"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name ActorProvider
type ActorProvider interface {
	EditActor(actor *models.Actor) error
	AddActor(actor *models.Actor) error
	AddMoviesToActor(actorID int64, movies []int64) error
	GetActors() ([]*models.ActorListing, error)
	DeleteActor(id int64) error
}

// @Summary Edit actor's data
// @Security ApiKeyAuth
// @Description Edit actor's data.
// @Tags Actors
// @Accept json
// @Produce json
// @Param input body models.editActor true "Actor object to be edited"
// @Success 200 {string} string "Successfully edited an actor"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /edit/actor [post]
func (h *Handler) editActor(w http.ResponseWriter, r *http.Request) {
	const op = "handler.editActor"

	log := h.log.With(slog.String("op", op))

	actor := &models.Actor{}
	err := json.NewDecoder(r.Body).Decode(actor)
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	log.Info("request body decoded")

	err = h.actorProvider.EditActor(actor)
	if err != nil {
		log.Error("failed to edit an actor", sl.Err(err))
		http.Error(w, "failed to edit an actor", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully edited an actor"))
}

// @Summary Get list of actors
// @Description Retrieves a list of actors.
// @Tags Actors
// @Accept json
// @Produce json
// @Success 200 {array} models.getActor "Successfully fetched actors"
// @Failure 500 {string} string "Internal server error"
// @Router /get/actors [get]
func (h *Handler) getActors(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getActors"

	log := h.log.With(slog.String("op", op))

	actors, err := h.actorProvider.GetActors()
	if err != nil {
		log.Error("failed to fetch actors", sl.Err(err))
		http.Error(w, "Failed to fetch actors", http.StatusInternalServerError)
		return
	}

	log.Info("fetched actors")

	actorJSON, err := json.Marshal(actors)
	if err != nil {
		log.Error("failed to marshal JSON", sl.Err(err))
		http.Error(w, "failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(actorJSON)
}

// @Summary Delete actor by ID
// @Security ApiKeyAuth
// @Description Deletes an actor by its ID.
// @Tags Actors
// @Accept json
// @Produce json
// @Param id query int true "Actor ID to be deleted"
// @Success 200 {string} string "Successfully deleted an actor"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /delete/actor [delete]
func (h *Handler) deleteActor(w http.ResponseWriter, r *http.Request) {
	const op = "handler.deleteActor"

	log := h.log.With(slog.String("op", op))

	actorIDStr := r.URL.Query().Get("id")
	actorID, err := strconv.ParseInt(actorIDStr, 10, 64)
	if err != nil {
		log.Error("invalid actor ID", sl.Err(err))
		http.Error(w, "invalid actor ID", http.StatusBadRequest)
		return
	}

	log.Info("parsed ID")

	err = h.actorProvider.DeleteActor(actorID)
	if err != nil {
		log.Error("failed to delete an actor", sl.Err(err))
		http.Error(w, "failed to delete an actor", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted an actor"))
}

// @Summary Add new actor
// @Security ApiKeyAuth
// @Description Adds a new actor using the provided actor object.
// @Tags Actors
// @Accept json
// @Produce json
// @Param input body models.addActor true "Actor object to be added"
// @Success 201 {string} string "Successfully added an actor"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /add/actor [post]
func (h *Handler) addActor(w http.ResponseWriter, r *http.Request) {
	const op = "handler.addActor"

	log := h.log.With(slog.String("op", op))

	actor := &models.Actor{}
	err := json.NewDecoder(r.Body).Decode(actor)
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	if actor.Name == "" {
		log.Error("actor name is empty", sl.Err(err))
		http.Error(w, "actor name is empty", http.StatusBadRequest)
		return
	}
	if actor.Sex == "" {
		log.Error("actor sex is empty", sl.Err(err))
		http.Error(w, "actor sex is empty", http.StatusBadRequest)
		return
	}
	if actor.Birthday.IsZero() {
		log.Error("actor birthday is empty", sl.Err(err))
		http.Error(w, "actor birthday is empty", http.StatusBadRequest)
		return
	}

	log.Info("request body decoded")

	err = h.actorProvider.AddActor(actor)
	if err != nil {
		log.Error("failed to add an actor", sl.Err(err))
		http.Error(w, "failed to add an actor", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Successfully added an actor"))
}

// @Summary Add movies to actor
// @Security ApiKeyAuth
// @Description Adds movies to an actor based on the provided actor ID and movie IDs.
// @Tags Actors
// @Accept json
// @Produce json
// @Param input body models.MoviesTo true "Actor ID and movie IDs to be added"
// @Success 200 {string} string "Successfully added movie(s) to actor"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /actor/add/movies [post]
func (h *Handler) addMoviesToActor(w http.ResponseWriter, r *http.Request) {
	const op = "handler.addMoviesToActor"

	log := h.log.With(slog.String("op", op))

	mtoa := &models.MoviesTo{}

	err := json.NewDecoder(r.Body).Decode(mtoa)
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	log.Info("request body decoded")

	err = h.actorProvider.AddMoviesToActor(mtoa.ActorID, mtoa.Movies)
	if err != nil {
		log.Error("failed to add movie(s) to actor", sl.Err(err))
		http.Error(w, "failed to add movie(s) to actor", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully added movie(s) to actor"))
}
