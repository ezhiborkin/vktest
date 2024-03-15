package models

import "time"

type Actor struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name,omitempty"`
	Sex       string    `json:"sex,omitempty"`
	Birthday  time.Time `json:"birthday,omitempty"`
	MoviesID  []int     `json:"movies_id,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

type ActorListing struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name,omitempty"`
	Sex      string    `json:"sex,omitempty"`
	Birthday time.Time `json:"birthday,omitempty"`
	Movies   []string  `json:"movies,omitempty"`
	//DeletedAt time.Time `json:"deleted_at,omitempty"`
}
