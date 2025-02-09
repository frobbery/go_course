package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrDateBusy = errors.New("date already has event")

	ErrEventNotFound = errors.New("no event found to update")
)

type Event struct {
	ID int64

	Name string

	DateTime time.Time

	EndDateTime time.Time

	Description string

	UserID int64

	SendBefore time.Time
}

type Notification struct {
	EventID int64

	Name string

	EventDate time.Time

	UserToSendID int64
}

type CalendarRepo interface {
	CreateEvent(ctx context.Context, event Event) (eventID int64, err error)

	UpdateEvent(ctx context.Context, eventID int64, event Event) (err error)

	DeleteEvent(ctx context.Context, eventID int64) (err error)

	EventsForDay(ctx context.Context, day time.Time) (events []Event, err error)

	EventsForWeek(ctx context.Context, firstDayOfWeek time.Time) (events []Event, err error)

	EventsForMonth(ctx context.Context, firstDayOfMonth time.Time) (events []Event, err error)
}
