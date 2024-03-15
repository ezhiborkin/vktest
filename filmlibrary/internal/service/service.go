package service

import (
	"filmlibrary/internal/domain/models"
	"filmlibrary/internal/lib/logger/sl"
	"fmt"
	"log/slog"
)

type Service struct {
	log          *slog.Logger
	actorStorage ActorStorage
	movieStorage MovieStorage
}

type ActorStorage interface {
	EditActorStorage(actor *models.Actor) error
	AddActorStorage(actor *models.Actor) error
	DeleteActorStorage(id int64) error
	GetActorsStorage() ([]*models.ActorListing, error)
}

type MovieStorage interface {
	EditMovieStorage(movie *models.Movie) error
	AddMovieStorage(movie *models.Movie) error
	DeleteMovieStorage(id int64) error
}

func New(log *slog.Logger, actorStorage ActorStorage, movieStorage MovieStorage) *Service {
	return &Service{log: log, actorStorage: actorStorage, movieStorage: movieStorage}
}

func (s *Service) AddActor(actor *models.Actor) error {
	const op = "service.AddActor"

	err := s.actorStorage.AddActorStorage(actor)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) GetActorsStorage() ([]*models.ActorListing, error) {
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
