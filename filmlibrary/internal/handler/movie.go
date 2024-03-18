package handler

import (
	"encoding/json"
	"errors"
	"filmlibrary/internal/domain/models"
	"filmlibrary/internal/lib/logger/sl"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name MovieProvider
type MovieProvider interface {
	GetMovie(input string) ([]*models.MovieListing, error)
	GetMoviesSorted(sortBy string, sortDirection string) ([]*models.MovieListing, error)
	EditMovie(movie *models.Movie) error
	AddMovie(movie *models.Movie) error
	AddActorsToMovie(movieID int64, actors []int64) error
	DeleteMovie(id int64) error
}

// @Summary Add movie
// @Security ApiKeyAuth
// @Description Adds a movie using the provided movie object.
// @Tags Movies
// @Accept json
// @Produce json
// @Param input body models.addMovie true "Movie object to be added"
// @Success 201 {string} string "Successfully added a movie"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /add/movie [post]
func (h *Handler) addMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.addMovie"

	log := h.log.With(slog.String("op", op))

	movie := &models.Movie{}
	err := json.NewDecoder(r.Body).Decode(movie)
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

	err = h.movieProvider.AddMovie(movie)
	if err != nil {
		log.Error("failed to add a movie", sl.Err(err))
		http.Error(w, "failed to add a movie", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Successfully added a movie"))
}

// @Summary Add actors to movie
// @Security ApiKeyAuth
// @Description Adds actors to a movie based on the provided movie ID and actor IDs.
// @Tags Movies
// @Accept json
// @Produce json
// @Param input body models.ActorsTo true "Movie ID and actor IDs to be added"
// @Success 200 {string} string "Successfully added actor(s) to movie"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /movie/add/actors [post]
func (h *Handler) addActorsToMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.addActorsToMovie"

	log := h.log.With(slog.String("op", op))

	atom := &models.ActorsTo{}

	err := json.NewDecoder(r.Body).Decode(atom)
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

	err = h.movieProvider.AddActorsToMovie(atom.MovieID, atom.Actors)
	if err != nil {
		log.Error("failed to add actors to a movie", sl.Err(err))
		http.Error(w, "failed to add actors to a movie", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully added actors to a movie"))
}

// @Summary Get movies sorted
// @Description Retrieves movies sorted by the provided criteria.
// @Tags Movies
// @Accept json
// @Produce json
// @Param sortBy query string false "Field to sort by"
// @Param sortDir query string false "Sort direction: asc or desc"
// @Success 200 {array} models.MovieListing "Sorted movies"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /get/movies [get]
func (h *Handler) getMoviesSorted(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getMoviesSorted"

	log := h.log.With(slog.String("op", op))

	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")

	movies, err := h.movieProvider.GetMoviesSorted(sortBy, sortDir)
	if err != nil {
		log.Error("failed to fetch movies", sl.Err(err))
		http.Error(w, "failed to fetch movies", http.StatusBadRequest)
		return
	}

	moviesJSON, err := json.Marshal(movies)
	if err != nil {
		log.Error("failed to marshal JSON", sl.Err(err))
		http.Error(w, "failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sorted by: " + sortBy + " " + sortDir))
	w.Write(moviesJSON)
}

// @Summary Get movie information
// @Description Get movie information based on substring of a title or an actor's name
// @Tags Movies
// @Accept json
// @Produce json
// @Param input query string true "Input to search for a movie"
// @Success 200 {array} models.MovieListing
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /find/movie [post]
func (h *Handler) getMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getMovie"

	log := h.log.With(slog.String("op", op))

	input := r.URL.Query().Get("input")
	movies, err := h.movieProvider.GetMovie(input)
	if err != nil {
		log.Error("failed to find a movie", sl.Err(err))
		http.Error(w, "failed to find a movie", http.StatusBadRequest)
		return
	}

	movieJSON, err := json.Marshal(movies)
	if err != nil {
		log.Error("failed to marshal JSON", sl.Err(err))
		http.Error(w, "failed to marshal JSON", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(movieJSON)
}

// @Summary Edit movie
// @Security ApiKeyAuth
// @Description Edit movie information.
// @Tags Movies
// @Accept json
// @Produce json
// @Param input body models.editMovie true "Movie object to be edited"
// @Success 200 {string} string "Successfully edited a movie"
// @Failure 400 {string} string "Bad request"
// @Failure 501 {string} string "Not Implemented"
// @Failure 500 {string} string "Internal server error"
// @Router /edit/movie [post]
func (h *Handler) editMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.editActor"

	log := h.log.With(slog.String("op", op))

	movie := &models.Movie{}
	err := json.NewDecoder(r.Body).Decode(movie)
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

	err = h.movieProvider.EditMovie(movie)
	if err != nil {
		log.Error("failed to edit a movie", sl.Err(err))
		http.Error(w, "failed to edit a movie", http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully edited a movie"))
}

// @Summary Delete movie
// @Security ApiKeyAuth
// @Description Delete a movie by its ID.
// @Tags Movies
// @Accept json
// @Produce json
// @Param id query int true "Movie ID to be deleted"
// @Success 200 {string} string "Successfully deleted a movie"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /delete/movie [delete]
func (h *Handler) deleteMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.deleteMovie"

	log := h.log.With(slog.String("op", op))

	movieIDStr := r.URL.Query().Get("id")
	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		log.Error("invalid movie ID", sl.Err(err))
		http.Error(w, "invalid movie ID", http.StatusBadRequest)
		return
	}

	log.Info("parsed movie ID")

	err = h.movieProvider.DeleteMovie(movieID)
	if err != nil {
		log.Error("failed to delete a movie", sl.Err(err))
		http.Error(w, "failed to delete a movie", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted a movie"))
}
