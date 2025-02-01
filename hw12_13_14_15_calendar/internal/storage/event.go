package storage

import (
	"context"
	"time"
)

type Event struct {
	Id int64
	Name string
	DateTime time.Time
	EndDateTime time.Time
	Description string
	UserId int64
	SendBefore time.Time
}

type Notification struct {
	EventId int64
	Name string
	EventDate time.Time
	UserToSendId int64
}

type CalendarRepo interface {
	CreateEvent(ctx context.Context, event Event) (eventId int64, err error)
	UpdateEvent(ctx context.Context, eventId int64, event Event) (err error)
	DeleteEvent(ctx context.Context, eventId int64) (err error)
	EventsForDay(ctx context.Context, day time.Time) (events []Event, err error)
	EventsForWeek(ctx context.Context, firstDayOfWeek time.Time) (events []Event, err error)
	EventsForMonth(ctx context.Context, firstDayOfMonth time.Time) (events []Event,err error)
}

type Storage interface {
	Connect(ctx context.Context, dsn string) error
	Migrate(ctx context.Context, migrate string) error
	Close() error
	CalendarRepo
}
