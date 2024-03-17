package handler

import (
	"context"
	"encoding/json"
	"errors"
	"filmlibrary/internal/domain/models"
	"filmlibrary/internal/lib/logger/sl"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	log           *slog.Logger
	serviceCooker ServiceCooker
}

type ServiceCooker interface {
	EditActor(actor *models.Actor) error
	AddActor(actor *models.Actor) error
	AddMoviesToActor(actorID int64, movies []int64) error
	GetActors() ([]*models.ActorListing, error)
	GetMovie(input string) ([]*models.MovieListing, error)
	GetMoviesSorted(sortBy string, sortDirection string) ([]*models.MovieListing, error)
	EditMovie(movie *models.Movie) error
	AddMovie(movie *models.Movie) error
	AddActorsToMovie(movieID int64, actors []int64) error
	DeleteActor(id int64) error
	DeleteMovie(id int64) error
	CreateUser(email, role, password string) error
	LoginUser(email, password string) (string, error)
}

func New(log *slog.Logger, serviceCooker ServiceCooker) *Handler {
	return &Handler{log: log, serviceCooker: serviceCooker}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/edit/actor", authMiddleware(onlyPostMiddleware(h.editActor)))
	mux.HandleFunc("/edit/movie", authMiddleware(onlyPostMiddleware(h.editMovie)))

	mux.HandleFunc("/create/actor", authMiddleware(onlyPostMiddleware(h.addActor)))
	mux.HandleFunc("/create/movie", authMiddleware(onlyPostMiddleware(h.addMovie)))
	mux.HandleFunc("/create/user", authMiddleware(onlyPostMiddleware(h.createUser)))

	mux.HandleFunc("/delete/actor", authMiddleware(onlyDeleteMiddleware(h.deleteActor)))
	mux.HandleFunc("/delete/movie", authMiddleware(onlyDeleteMiddleware(h.deleteMovie)))

	mux.HandleFunc("/actor/add/movies", authMiddleware(onlyPostMiddleware(h.addMoviesToActor)))
	mux.HandleFunc("/movie/add/actors", authMiddleware(onlyPostMiddleware(h.addActorsToMovie)))

	mux.HandleFunc("/get/actors", onlyGetMiddleware(h.getActors))
	mux.HandleFunc("/get/movies", onlyGetMiddleware(h.getMoviesSorted))

	mux.HandleFunc("/find/movie", onlyPostMiddleware(h.getMovie)) //

	mux.HandleFunc("/login", onlyPostMiddleware(h.loginUser))

	return mux
}

func (h *Handler) editActor(w http.ResponseWriter, r *http.Request) {
	const op = "handler.editActor"

	log := h.log.With(slog.String("op", op))

	actor := &models.Actor{}
	err := json.NewDecoder(r.Body).Decode(actor)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))
			http.Error(w, "empty request", http.StatusBadRequest)
			return
		}
		log.Error("failed to decode request body", sl.Err(err))
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	log.Info("request body decoded")

	err = h.serviceCooker.EditActor(actor)
	if err != nil {
		log.Error("failed to edit an actor", sl.Err(err))
		http.Error(w, "failed to edit an actor", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully edited an actor"))
}

func (h *Handler) getActors(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getActors"

	log := h.log.With(slog.String("op", op))

	actors, err := h.serviceCooker.GetActors()
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

	err = h.serviceCooker.DeleteActor(actorID)
	if err != nil {
		log.Error("failed to delete an actor", sl.Err(err))
		http.Error(w, "failed to delete an actor", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted an actor"))
}

func (h *Handler) addActor(w http.ResponseWriter, r *http.Request) {
	const op = "handler.addActor"

	log := h.log.With(slog.String("op", op))

	actor := &models.Actor{}
	err := json.NewDecoder(r.Body).Decode(actor)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))
			http.Error(w, "empty request", http.StatusBadRequest)
			return
		}
		log.Error("failed to decode request body", sl.Err(err))
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	log.Info("request body decoded")

	err = h.serviceCooker.AddActor(actor)
	if err != nil {
		log.Error("failed to add an actor", sl.Err(err))
		http.Error(w, "failed to add an actor", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Successfully added an actor"))
}

func (h *Handler) addMoviesToActor(w http.ResponseWriter, r *http.Request) {
	const op = "handler.addMoviesToActor"

	log := h.log.With(slog.String("op", op))

	var req struct {
		ActorID int64   `json:"id"`
		Movies  []int64 `json:"movies_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))
			http.Error(w, "empty request", http.StatusBadRequest)
			return
		}
		log.Error("failed to decode request body", sl.Err(err))
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	log.Info("request body decoded")

	err = h.serviceCooker.AddMoviesToActor(req.ActorID, req.Movies)
	if err != nil {
		log.Error("failed to add movie(s) to actor", sl.Err(err))
		http.Error(w, "failed to add movie(s) to actor", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully added movie(s) to actor"))
}

func (h *Handler) addMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.addMovie"

	log := h.log.With(slog.String("op", op))

	movie := &models.Movie{}
	err := json.NewDecoder(r.Body).Decode(movie)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))
			http.Error(w, "empty request", http.StatusBadRequest)
			return
		}
		log.Error("failed to decode request body", sl.Err(err))
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	log.Info("request body decoded")

	err = h.serviceCooker.AddMovie(movie)
	if err != nil {
		log.Error("failed to add a movie", sl.Err(err))
		http.Error(w, "failed to add a movie", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Successfully added a movie"))
}

func (h *Handler) addActorsToMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.addActorsToMovie"

	log := h.log.With(slog.String("op", op))

	var req struct {
		MovieID int64   `json:"id"`
		Actors  []int64 `json:"actors_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))
			http.Error(w, "empty request", http.StatusBadRequest)
			return
		}
		log.Error("failed to decode request body", sl.Err(err))
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	log.Info("request body decoded")

	err = h.serviceCooker.AddActorsToMovie(req.MovieID, req.Actors)
	if err != nil {
		log.Error("failed to add actors to a movie", sl.Err(err))
		http.Error(w, "failed to add actors to a movie", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully added actors to a movie"))
}

func (h *Handler) getMoviesSorted(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getMoviesSorted"

	log := h.log.With(slog.String("op", op))

	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")

	movies, err := h.serviceCooker.GetMoviesSorted(sortBy, sortDir)
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

func (h *Handler) getMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getMovie"

	log := h.log.With(slog.String("op", op))

	input := r.URL.Query().Get("input")
	movies, err := h.serviceCooker.GetMovie(input)
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
	w.Write([]byte("Found possible movies: \n"))
	w.Write(movieJSON)
}

func (h *Handler) editMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.editActor"

	log := h.log.With(slog.String("op", op))

	movie := &models.Movie{}
	err := json.NewDecoder(r.Body).Decode(movie)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))
			http.Error(w, "empty request", http.StatusBadRequest)
			return
		}
		log.Error("failed to decode request body", sl.Err(err))
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	log.Info("request body decoded")

	err = h.serviceCooker.EditMovie(movie)
	if err != nil {
		log.Error("failed to edit a movie", sl.Err(err))
		http.Error(w, "failed to edit a movie", http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully edited a movie"))
}

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

	err = h.serviceCooker.DeleteMovie(movieID)
	if err != nil {
		log.Error("failed to delete a movie", sl.Err(err))
		http.Error(w, "failed to delete a movie", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted a movie"))
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	const op = "handler.createUser"

	log := h.log.With(slog.String("op", op))

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))
			http.Error(w, "empty request", http.StatusBadRequest)
			return
		}
		log.Error("failed to decode request body", sl.Err(err))
		http.Error(w, "failed to decode request", http.StatusBadRequest)
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

	err = h.serviceCooker.CreateUser(user.Email, user.Role, user.Password)
	if err != nil {
		log.Error("failed to create user", sl.Err(err))
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Created user with email - " + user.Email))
}

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	const op = "handler.loginUser"

	log := h.log.With(slog.String("op", op))

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))
			http.Error(w, "empty request", http.StatusBadRequest)
			return
		}
		log.Error("failed to decode request body", sl.Err(err))
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	log.Info("request body decoded")

	tokenString, err := h.serviceCooker.LoginUser(user.Email, user.Password)
	if err != nil {
		log.Error("failed to login user", sl.Err(err))
		http.Error(w, "failed to login user", http.StatusBadRequest)
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully logged in, your token: " + tokenString))
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

func authMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			fmt.Println("Authorization header missing")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authorization header missing"))
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})

		if err != nil {
			fmt.Println("Error parsing token:", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			role, ok := claims["Role"].(string)
			if !ok {
				fmt.Println("Role claim not found or not a string")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Role claim not found or not a string"))
				return
			}

			if role != "admin" {
				fmt.Println("Wrong role")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Wrong role"))
				return
			}

			ctx := context.WithValue(r.Context(), "role", role)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			fmt.Println("Invalid token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid token"))
		}
	}
}
