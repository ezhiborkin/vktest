package handler

import (
	"context"
	_ "filmlibrary/docs"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log/slog"
	"net/http"
	"strings"
)

type Handler struct {
	log           *slog.Logger
	userProvider  UserProvider
	actorProvider ActorProvider
	movieProvider MovieProvider
	authProvider  AuthProvider
}

func New(log *slog.Logger,
	userProvider UserProvider,
	actorProvider ActorProvider,
	movieProvider MovieProvider,
	authProvider AuthProvider,
) *Handler {
	return &Handler{
		log:           log,
		userProvider:  userProvider,
		actorProvider: actorProvider,
		movieProvider: movieProvider,
		authProvider:  authProvider,
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("/edit/actor", authMiddleware(onlyPostMiddleware(h.editActor)))
	mux.HandleFunc("/edit/movie", authMiddleware(onlyPostMiddleware(h.editMovie)))

	mux.HandleFunc("/create/actor", authMiddleware(onlyPostMiddleware(h.addActor)))
	mux.HandleFunc("/create/movie", authMiddleware(onlyPostMiddleware(h.addMovie)))

	mux.HandleFunc("/delete/actor", authMiddleware(onlyDeleteMiddleware(h.deleteActor)))
	mux.HandleFunc("/delete/movie", authMiddleware(onlyDeleteMiddleware(h.deleteMovie)))

	mux.HandleFunc("/actor/add/movies", authMiddleware(onlyPostMiddleware(h.addMoviesToActor)))
	mux.HandleFunc("/movie/add/actors", authMiddleware(onlyPostMiddleware(h.addActorsToMovie)))

	mux.HandleFunc("/get/actors", onlyGetMiddleware(h.getActors))
	mux.HandleFunc("/get/movies", onlyGetMiddleware(h.getMoviesSorted))

	mux.HandleFunc("/find/movie", onlyPostMiddleware(h.getMovie))

	mux.HandleFunc("/login", onlyPostMiddleware(h.loginUser))
	mux.HandleFunc("/create/user", onlyPostMiddleware(h.createUser))

	return mux
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
