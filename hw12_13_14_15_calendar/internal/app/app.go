package app

import (
	"context"
	"github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
}

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

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
