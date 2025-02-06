package logger

import (
	"bytes"
	"io"
	"log"
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("test info log", func(t *testing.T) {
		buf := bytes.Buffer{}
		logger := log.New(&buf, "", 0)

		logg := New("INFO")
		logg.logger = logger

		logg.Info("info")
		logLine, err := buf.ReadString(byte('\n'))
		require.Nil(t, err)
		require.Equal(t, logLine, "[INFO]info\n")

		logg.Error("error")
		logLine, err = buf.ReadString(byte('\n'))
		require.Nil(t, err)
		require.Equal(t, logLine, "[ERROR]error\n")

		logg.Debug("debug")
		logLine, err = buf.ReadString(byte('\n'))
		require.Equal(t, err, io.EOF)
		require.Equal(t, logLine, "")
	})

	t.Run("test error log", func(t *testing.T) {
		buf := bytes.Buffer{}
		logger := log.New(&buf, "", 0)

		logg := New("ERROR")
		logg.logger = logger

		logg.Info("info")
		logLine, err := buf.ReadString(byte('\n'))
		require.Equal(t, err, io.EOF)
		require.Equal(t, logLine, "")

		logg.Error("error")
		logLine, err = buf.ReadString(byte('\n'))
		require.Nil(t, err)
		require.Equal(t, logLine, "[ERROR]error\n")

		logg.Debug("debug")
		logLine, err = buf.ReadString(byte('\n'))
		require.Equal(t, err, io.EOF)
		require.Equal(t, logLine, "")
	})

	t.Run("test debug log", func(t *testing.T) {
		buf := bytes.Buffer{}
		logger := log.New(&buf, "", 0)

		logg := New("DEBUG")
		logg.logger = logger

		logg.Info("info")
		logLine, err := buf.ReadString(byte('\n'))
		require.Nil(t, err)
		require.Equal(t, logLine, "[INFO]info\n")

		logg.Error("error")
		logLine, err = buf.ReadString(byte('\n'))
		require.Nil(t, err)
		require.Equal(t, logLine, "[ERROR]error\n")

		logg.Debug("debug")
		logLine, err = buf.ReadString(byte('\n'))
		require.Nil(t, err)
		require.Equal(t, logLine, "[DEBUG]debug\n")
	})
}
