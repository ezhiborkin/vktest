package models

import "time"

type Movie struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	ReleaseDate time.Time `json:"release_date,omitempty"`
	Rating      *float64  `json:"rating,omitempty"`
	ActorsID    []int     `json:"actors_id,omitempty"`
	DeletedAt   time.Time `json:"deleted_at,omitempty"`
}

type MovieListing struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	ReleaseDate time.Time `json:"release_date,omitempty"`
	Rating      *float64  `json:"rating,omitempty"`
	Actors      []string  `json:"actors_id,omitempty"`
}

type MoviesTo struct {
	ActorID int64   `json:"id" binding:"required" example:"1"`
	Movies  []int64 `json:"movies_id,omitempty" binding:"required" example:"123"`
}

type editMovie struct {
	ID          int64     `json:"id" binding:"required" example:"1"`
	Title       string    `json:"title,omitempty" example:"The Shawshank Redemption"`
	Description string    `json:"description,omitempty" example:"Two"`
	ReleaseDate time.Time `json:"release_date,omitempty" example:"1994-10-14" format:"date"`
	Rating      *float64  `json:"rating,omitempty" example:"9.3"`
}

type addMovie struct {
	Title       string    `json:"title,omitempty" binding:"required" example:"The Shawshank Redemption"`
	Description string    `json:"description,omitempty" binding:"required" example:"Two"`
	ReleaseDate time.Time `json:"release_date,omitempty" binding:"required" example:"1994-10-14" format:"date"`
	Rating      *float64  `json:"rating,omitempty" binding:"required" example:"9.3"`
	ActorsID    []int     `json:"actors_id,omitempty" binding:"required"`
}
