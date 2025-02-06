package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/app"
	"github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/storage"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type sqlStorage struct { 
	db *sql.DB
}

func New() app.Storage {
	return &sqlStorage{}
}

func (s *sqlStorage) Connect(ctx context.Context, dsn string) (err error) {
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}
	return s.db.PingContext(ctx)
}

///go:embed migrations/*.sql
//var embedMigrations embed.FS

func (s *sqlStorage) Migrate(ctx context.Context, migrate string) (err error) {
	//	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}

	if err := goose.Up(s.db, migrate); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}

func (s *sqlStorage) Close() error {
	return s.db.Close()
}

func (s *sqlStorage) CreateEvent(ctx context.Context, event storage.Event) (eventId int64, err error) {
	eventsForDay, err := s.EventsForDay(ctx, event.DateTime)
	if err != nil {
		return -1, err
	}
	if len(eventsForDay) != 0 {
		return -1, storage.ErrDateBusy
	}
	query := `insert into event(name, date_time, end_date_time, description, user_id, send_before) values($1, $2, $3, $4, $5, $6)`
	result, err := s.db.ExecContext(ctx, query, event.Name, event.DateTime, event.EndDateTime, event.Description, event.UserId, event.SendBefore)
	if err != nil {
 		return 0, nil
	}
	eventId, err = result.LastInsertId() 
	return eventId, err
}

func (s *sqlStorage) UpdateEvent(ctx context.Context, eventId int64, event storage.Event) (err error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id FROM event
		where id = $1
	`, eventId)
	if err != nil {
		return fmt.Errorf("cannot select: %w", err)
	}
	defer rows.Close()
	if !rows.Next() {
		return storage.ErrEventNotFound
	}
	query := `update event set name = $2, date_time = $3, end_date_time = $4, description = $5, user_id = $6, send_before = $7 where id = $1`
	_, err = s.db.ExecContext(ctx, query, eventId, event.Name, event.DateTime, event.EndDateTime, event.Description, event.UserId, event.SendBefore)
	return err
}

func (s *sqlStorage) DeleteEvent(ctx context.Context, eventId int64) (err error) {
	query := `delete from event where id = $1`
	_, err = s.db.ExecContext(ctx, query, eventId)
	return err
}

func(s *sqlStorage) getEventsBetweenTwoDates(ctx context.Context, startDate time.Time, endDate time.Time) (events []storage.Event, err error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, date_time, end_date_time, description, user_id, send_before FROM event
		where tstzrange($1, $2) @> date_time
	`, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("cannot select: %w", err)
	}
	defer rows.Close()

	events = make([]storage.Event, 0)

	for rows.Next() {
		var e storage.Event
		var description sql.NullString
		var sendBefore sql.NullTime

		if err := rows.Scan(
			&e.Id,
			&e.Name,
			&e.DateTime,
			&e.EndDateTime,
			&description,
			&e.UserId,
			&sendBefore,
		); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}

		if description.Valid {
			e.Description = e.Description
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (s *sqlStorage) EventsForDay(ctx context.Context, day time.Time) (events []storage.Event, err error) {
	return s.getEventsBetweenTwoDates(ctx, day, day.Add(time.Duration(time.Hour * 24)));
}

func (s *sqlStorage) EventsForWeek(ctx context.Context, day time.Time) (events []storage.Event, err error) {
	return s.getEventsBetweenTwoDates(ctx, day, day.Add(time.Duration(time.Hour * 24 * 7)));
}

func (s *sqlStorage) EventsForMonth(ctx context.Context, day time.Time) (events []storage.Event, err error) {
	return s.getEventsBetweenTwoDates(ctx, day, day.Add(time.Duration(time.Hour * 24 * 30)));
}
