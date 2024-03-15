package service

import (
	"filmlibrary/internal/domain/models"
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
}

type MovieStorage interface {
	EditMovieStorage(movie *models.Movie) error
	AddMovieStorage(movie *models.Movie) error
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

func (s *Service) EditActor(actor *models.Actor) error {
	const op = "service.EditActor"

	err := s.actorStorage.EditActorStorage(actor)
	if err != nil {
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
