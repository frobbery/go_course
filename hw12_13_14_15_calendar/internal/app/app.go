package app

import (
	"context"

	//nolint:depguard
	"github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/storage"
)

type App struct{}

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string)
}

type Storage interface {
	Connect(ctx context.Context, dsn string) error
	Migrate(ctx context.Context, migrate string) error
	Close() error
	storage.CalendarRepo
}

func New(_ Logger, _ Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(_ context.Context, _, _ string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
