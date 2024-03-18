package service

import (
	"filmlibrary/internal/domain/models"
	"fmt"
)

type MovieStorage interface {
	EditMovieStorage(movie *models.Movie) error
	AddMovieStorage(movie *models.Movie) error
	DeleteMovieStorage(id int64) error
	GetMoviesSortedStorage(sortBy string, sortDirection string) ([]*models.MovieListing, error)
	AddActorsToMovieStorage(movieID int64, actors []int64) error
	GetMovieStorage(searchTerm string) ([]*models.MovieListing, error)
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

	movies, err := s.movieStorage.GetMovieStorage(input)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movies, nil
}

func (s *Service) EditMovie(movie *models.Movie) error {
	const op = "service.EditMovie"

	err := s.movieStorage.EditMovieStorage(movie)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteMovie(id int64) error {
	const op = "service.DeleteMovie"

	err := s.movieStorage.DeleteMovieStorage(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
