package service

import (
	"filmlibrary/internal/domain/models"
	"fmt"
)

type ActorStorage interface {
	EditActorStorage(actor *models.Actor) error
	AddActorStorage(actor *models.Actor) error
	DeleteActorStorage(id int64) error
	GetActorsStorage() ([]*models.ActorListing, error)
	AddMoviesToActorStorage(actorID int64, movies []int64) error
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

	err := s.actorStorage.DeleteActorStorage(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
