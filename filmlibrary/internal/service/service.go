package service

import (
	"log/slog"
)

type Service struct {
	log          *slog.Logger
	actorStorage ActorStorage
	movieStorage MovieStorage
	userStorage  UserStorage
}

func New(log *slog.Logger, actorStorage ActorStorage, movieStorage MovieStorage, userStorage UserStorage) *Service {
	return &Service{log: log, actorStorage: actorStorage, movieStorage: movieStorage, userStorage: userStorage}
}
