package service

import (
	"errors"
	"filmlibrary/internal/domain/models"
	"filmlibrary/internal/lib/logger/sl"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Service struct {
	log          *slog.Logger
	actorStorage ActorStorage
	movieStorage MovieStorage
	userStorage  UserStorage
}

type ActorStorage interface {
	EditActorStorage(actor *models.Actor) error
	AddActorStorage(actor *models.Actor) error
	DeleteActorStorage(id int64) error
	GetActorsStorage() ([]*models.ActorListing, error)
	AddMoviesToActorStorage(actorID int64, movies []int64) error
}

type MovieStorage interface {
	EditMovieStorage(movie *models.Movie) error
	AddMovieStorage(movie *models.Movie) error
	DeleteMovieStorage(id int64) error
	GetMoviesSortedStorage(sortBy string, sortDirection string) ([]*models.MovieListing, error)
	AddActorsToMovieStorage(movieID int64, actors []int64) error
	GetMovieStorage(searchTerm string) ([]*models.MovieListing, error)
}

type UserStorage interface {
	CreateUserStorage(email, role string, passHash []byte) error
	GetUserStorage(email string) (*models.User, error)
}

func New(log *slog.Logger, actorStorage ActorStorage, movieStorage MovieStorage, userStorage UserStorage) *Service {
	return &Service{log: log, actorStorage: actorStorage, movieStorage: movieStorage, userStorage: userStorage}
}

func (s *Service) AddActor(actor *models.Actor) error {
	const op = "service.AddActor"

	err := s.actorStorage.AddActorStorage(actor)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) AddMoviesToActor(actorID int64, movies []int64) error {
	const op = "service.AddMoviesToActor"

	err := s.actorStorage.AddMoviesToActorStorage(actorID, movies)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) GetActors() ([]*models.ActorListing, error) {
	const op = "service.GetActorsStorage"

	actors, err := s.actorStorage.GetActorsStorage()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return actors, nil
}

func (s *Service) EditActor(actor *models.Actor) error {
	const op = "service.EditActor"

	err := s.actorStorage.EditActorStorage(actor)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteActor(id int64) error {
	const op = "service.DeleteActor"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("deleting actor")

	err := s.actorStorage.DeleteActorStorage(id)
	if err != nil {
		s.log.Error("failed to delete an actor", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) AddMovie(movie *models.Movie) error {
	const op = "service.AddMovie"

	err := s.movieStorage.AddMovieStorage(movie)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) AddActorsToMovie(movieID int64, actors []int64) error {
	const op = "service.AddActorsToMovie"

	err := s.movieStorage.AddActorsToMovieStorage(movieID, actors)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) GetMoviesSorted(sortBy string, sortDirection string) ([]*models.MovieListing, error) {
	const op = "service.GetMoviesSorted"

	movies, err := s.movieStorage.GetMoviesSortedStorage(sortBy, sortDirection)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movies, nil
}

func (s *Service) GetMovie(input string) ([]*models.MovieListing, error) {
	const op = "service.GetMovie"

	fmt.Println("service")

	movies, err := s.movieStorage.GetMovieStorage(input)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movies, nil
}

func (s *Service) EditMovie(movie *models.Movie) error {
	const op = "service.EditMovie"

	fmt.Println(op)

	err := s.movieStorage.EditMovieStorage(movie)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteMovie(id int64) error {
	const op = "service.DeleteMovie"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("deleting movie")

	err := s.movieStorage.DeleteMovieStorage(id)
	if err != nil {
		s.log.Error("failed to delete a movie", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) CreateUser(email, role string, password string) error {
	const op = "service.CreateUser"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("creating user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to generate passHash", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.userStorage.CreateUserStorage(email, role, passHash)
	if err != nil {
		s.log.Error("failed to create a user", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) LoginUser(email, password string) (string, error) {
	const op = "service.LoginUser"

	expiresAt := time.Now().Add(time.Minute * 100).Unix()

	user, err := s.userStorage.GetUserStorage(email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	tk := &models.Token{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}
