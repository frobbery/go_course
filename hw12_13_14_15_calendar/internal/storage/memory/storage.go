package memorystorage

import (
	"context"
	"sync"
	"time"

	//nolint:depguard
	"github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/app"
	//nolint:depguard
	"github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/storage"
)

type inMemoryStorage struct {
	currentEventID int64
	eventsByID     map[int64]storage.Event
	mu             sync.RWMutex
}

func New() app.Storage {
	return &inMemoryStorage{}
}

func (s *inMemoryStorage) Connect(_ context.Context, _ string) (err error) {
	s.eventsByID = make(map[int64]storage.Event)
	return nil
}

func (s *inMemoryStorage) Migrate(_ context.Context, _ string) (err error) {
	return nil
}

func (s *inMemoryStorage) Close() error {
	for eventID := range s.eventsByID {
		delete(s.eventsByID, eventID)
	}
	return nil
}

func (s *inMemoryStorage) CreateEvent(ctx context.Context, event storage.Event) (eventID int64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	eventsForDay := s.getEventsBetweenTwoDates(ctx, event.DateTime, event.DateTime.Add(time.Hour*24))
	if len(eventsForDay) != 0 {
		return -1, storage.ErrDateBusy
	}
	s.currentEventID++
	event.ID = s.currentEventID
	s.eventsByID[s.currentEventID] = event
	return s.currentEventID, nil
}

func (s *inMemoryStorage) UpdateEvent(_ context.Context, eventID int64, event storage.Event) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, eventExists := s.eventsByID[eventID]
	if !eventExists {
		return storage.ErrEventNotFound
	}
	s.eventsByID[eventID] = event
	return nil
}

func (s *inMemoryStorage) DeleteEvent(_ context.Context, eventID int64) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.eventsByID, eventID)
	return nil
}

//nolint:lll
func (s *inMemoryStorage) getEventsBetweenTwoDates(_ context.Context, startDate time.Time, endDate time.Time) (events []storage.Event) {
	events = make([]storage.Event, 0)
	for eventID := range s.eventsByID {
		curDateTime := s.eventsByID[eventID].DateTime
		if !curDateTime.Before(startDate) && !curDateTime.After(endDate) {
			events = append(events, s.eventsByID[eventID])
		}
	}
	return events
}

func (s *inMemoryStorage) EventsForDay(ctx context.Context, day time.Time) (events []storage.Event, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getEventsBetweenTwoDates(ctx, day, day.Add(time.Hour*24)), nil
}

func (s *inMemoryStorage) EventsForWeek(ctx context.Context, day time.Time) (events []storage.Event, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getEventsBetweenTwoDates(ctx, day, day.Add(time.Hour*24*7)), nil
}

func (s *inMemoryStorage) EventsForMonth(ctx context.Context, day time.Time) (events []storage.Event, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getEventsBetweenTwoDates(ctx, day, day.Add(time.Hour*24*30)), nil
}
