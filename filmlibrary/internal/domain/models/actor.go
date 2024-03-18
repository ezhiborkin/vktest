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
}

type getActor struct {
	ID       int64     `json:"id" binding:"required" example:"1"`
	Name     string    `json:"name,omitempty"  example:"Vladimir Putin"`
	Sex      string    `json:"sex,omitempty"  example:"male"`
	Birthday time.Time `json:"birthday,omitempty"  example:"1952-10-07" format:"date"`
	Movies   []string  `json:"movies,omitempty" example:"['The Shawshank Redemption', 'The Godfather']"`
}

type editActor struct {
	ID       int64     `json:"id" binding:"required" example:"1"`
	Name     string    `json:"name,omitempty"  example:"Vladimir Putin"`
	Sex      string    `json:"sex,omitempty"  example:"male"`
	Birthday time.Time `json:"birthday,omitempty"  example:"1952-10-07" format:"date"`
}

type addActor struct {
	Name     string    `json:"name,omitempty" binding:"required" example:"Vladimir Putin"`
	Sex      string    `json:"sex,omitempty"  binding:"required" example:"male"`
	Birthday time.Time `json:"birthday,omitempty" binding:"required" example:"1952-10-07" format:"date"`
}

type ActorsTo struct {
	MovieID int64   `json:"id" binding:"required" example:"2"`
	Actors  []int64 `json:"actors_id,omitempty" binding:"required" example:"123"`
}
