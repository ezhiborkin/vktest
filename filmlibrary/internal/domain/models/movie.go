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
