package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/frobbery/go_course/hw12_13_14_15_calendar/internal/storage"
	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("test in memory storage", func(t *testing.T) {
		s := New()
		s.Connect(context.Background(), "")
		event := storage.Event{
			Id : 1,
			Name: "ny",
			DateTime: time.Date(2025, time.December, 31, 0, 0, 0, 0, time.Now().Location()),
			EndDateTime: time.Date(2025, time.December, 31, 23, 59, 59, 0, time.Now().Location()),
			Description: "new year",
			UserId: 1,
			SendBefore: time.Date(2025, time.December, 30, 12, 0, 0, 0, time.Now().Location()),
		}

		resId, err := s.CreateEvent(context.Background(), event)
		require.Nil(t, err)
		require.Equal(t, resId, int64(1))

		res, err := s.EventsForDay(context.Background(), time.Date(2025, time.December, 31, 0, 0, 0, 0, time.Now().Location()))
		require.Nil(t, err)
		require.Equal(t, res, []storage.Event{event})

		res, err = s.EventsForWeek(context.Background(), time.Date(2025, time.December, 30, 0, 0, 0, 0, time.Now().Location()))
		require.Nil(t, err)
		require.Equal(t, res, []storage.Event{event})

		res, err = s.EventsForMonth(context.Background(), time.Date(2025, time.December, 4, 0, 0, 0, 0, time.Now().Location()))
		require.Nil(t, err)
		require.Equal(t, res, []storage.Event{event})

		s.UpdateEvent(context.Background(), 1, storage.Event{
			Id : 1,
			Name: "ny",
			DateTime: time.Date(2025, time.December, 31, 0, 0, 0, 0, time.Now().Location()),
			EndDateTime: time.Date(2025, time.December, 31, 23, 59, 59, 0, time.Now().Location()),
			Description: "new description",
			UserId: 1,
			SendBefore: time.Date(2025, time.December, 30, 12, 0, 0, 0, time.Now().Location()),
		})
		res, err = s.EventsForDay(context.Background(), time.Date(2025, time.December, 31, 0, 0, 0, 0, time.Now().Location()))
		require.Nil(t, err)
		require.Equal(t, res[0].Description, "new description")

		err = s.DeleteEvent(context.Background(), 1)
		require.Nil(t, err)
		res, err = s.EventsForDay(context.Background(), time.Date(2025, time.December, 31, 0, 0, 0, 0, time.Now().Location()))
		require.Nil(t, err)
		require.Equal(t, len(res), 0)
	})
}
