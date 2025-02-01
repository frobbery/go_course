package memorystorage

import (
	"sync"
	"github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/storage"
	"context"
	"time"
)

type inMemoryStorage struct {
	currentEventId 	int64
	eventsById		map[int64]storage.Event
	mu				sync.RWMutex
}

func New() storage.Storage {
	return &inMemoryStorage{}
}

func (s *inMemoryStorage) Connect(ctx context.Context, dsn string) (err error) {
	s.eventsById = make(map[int64]storage.Event)
	return nil
}

func (s *inMemoryStorage) Migrate(ctx context.Context, migrate string) (err error) {
	return nil
}

func (s *inMemoryStorage) Close() error {
	for eventId := range s.eventsById {
		delete(s.eventsById, eventId)
	}
	return nil
}

func (s *inMemoryStorage) CreateEvent(ctx context.Context, event storage.Event) (eventId int64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentEventId++
	event.Id = s.currentEventId
	s.eventsById[s.currentEventId] = event
	return s.currentEventId, nil
}

func (s *inMemoryStorage) UpdateEvent(ctx context.Context, eventId int64, event storage.Event) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.eventsById[eventId] = event
	return nil
}

func (s *inMemoryStorage) DeleteEvent(ctx context.Context, eventId int64) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.eventsById, eventId)
	return nil
}

func(s *inMemoryStorage) getEventsBetweenTwoDates(_ context.Context, startDate time.Time, endDate time.Time) (events []storage.Event, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events = make([]storage.Event, 0)
	for eventId := range s.eventsById {
		curDateTime := s.eventsById[eventId].DateTime
		if !curDateTime.Before(startDate) && !curDateTime.After(endDate) {
			events = append(events, s.eventsById[eventId])
		}
	}
	return events, nil
}

func (s *inMemoryStorage) EventsForDay(ctx context.Context, day time.Time) (events []storage.Event, err error) {
	return s.getEventsBetweenTwoDates(ctx, day, day.Add(time.Duration(time.Hour * 24)));
}

func (s *inMemoryStorage) EventsForWeek(ctx context.Context, day time.Time) (events []storage.Event, err error) {
	return s.getEventsBetweenTwoDates(ctx, day, day.Add(time.Duration(time.Hour * 24 * 7)));
}

func (s *inMemoryStorage) EventsForMonth(ctx context.Context, day time.Time) (events []storage.Event, err error) {
	return s.getEventsBetweenTwoDates(ctx, day, day.Add(time.Duration(time.Hour * 24 * 30)));
}
