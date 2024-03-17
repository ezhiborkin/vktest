package handler

import (
	"context"
	"encoding/json"
	"filmlibrary/internal/domain/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
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

	mux.HandleFunc("/get/actors", authMiddleware(onlyGetMiddleware(h.getActors))) /////
	mux.HandleFunc("/get/movies", onlyGetMiddleware(h.getMoviesSorted))

	mux.HandleFunc("/actor/add/movies", onlyPostMiddleware(h.addMoviesToActor))
	mux.HandleFunc("/movie/add/actors", onlyPostMiddleware(h.addActorsToMovie))

	mux.HandleFunc("/find/movie", onlyPostMiddleware(h.getMovie))

	mux.HandleFunc("/create/user", onlyPostMiddleware(h.createUser))
	mux.HandleFunc("/login/user", onlyPostMiddleware(h.loginUser))

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

	actors, err := h.serviceCooker.GetActors()
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

func (h *Handler) addMoviesToActor(w http.ResponseWriter, r *http.Request) {
	const op = "handler.addMoviesToActor"

	var req struct {
		ActorID int64   `json:"id"`
		Movies  []int64 `json:"movies_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, op, http.StatusBadRequest)
	}

	err = h.serviceCooker.AddMoviesToActor(req.ActorID, req.Movies)
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

func (h *Handler) addActorsToMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.addActorsToMovie"

	var req struct {
		MovieID int64   `json:"id"`
		Actors  []int64 `json:"actors_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, op, http.StatusBadRequest)
	}

	err = h.serviceCooker.AddActorsToMovie(req.MovieID, req.Actors)
	if err != nil {
		http.Error(w, op, http.StatusNotImplemented)
	}

	response := map[string]interface{}{
		"kekys": 100,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) getMoviesSorted(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getMoviesSorted"

	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")

	movies, err := h.serviceCooker.GetMoviesSorted(sortBy, sortDir)
	if err != nil {
		http.Error(w, "Failed to fetch movies", http.StatusInternalServerError)
		fmt.Printf("%s: %w", op, err)
	}

	moviesJSON, err := json.Marshal(movies)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		fmt.Printf("%s: %v\n", op, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(moviesJSON)
}

func (h *Handler) getMovie(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getMovie"

	input := r.URL.Query().Get("input")

	movies, err := h.serviceCooker.GetMovie(input)
	if err != nil {
		http.Error(w, "Failed to find a movie", http.StatusInternalServerError)
		fmt.Printf("%s: %w", op, err)
	}

	movieJSON, err := json.Marshal(movies)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		fmt.Printf("%s: %v\n", op, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(movieJSON); err != nil {
		fmt.Printf("%s: %v\n", op, err)
	}
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

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	const op = "handler.createUser"

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, op, http.StatusBadRequest)
		return
	}

	err = h.serviceCooker.CreateUser(user.Email, user.Role, user.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		fmt.Printf("%s: %v\n", op, err)
		return
	}

	emailResponse := "Created user with email - " + user.Email

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "%s"}`, emailResponse)
}

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	const op = "handler.loginUser"

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, op, http.StatusBadRequest)
		return
	}

	fmt.Println("user email:", user.Email)

	tokenString, err := h.serviceCooker.LoginUser(user.Email, user.Password)
	if err != nil {
		http.Error(w, op, http.StatusInternalServerError)
		return
	}
	fmt.Println(tokenString)

	// Set the token in the response header
	w.Header().Set("Authorization", "Bearer "+tokenString)

	// Return a success status
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("your token: " + tokenString))
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
			// Now you can use the role in your middleware or pass it to the context
			fmt.Println("User role:", role)
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
