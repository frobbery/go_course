package main

import (
	"os"
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("Offset 0, limit 0", func(t *testing.T) {
		err := Copy("testdata/input1.txt", "testdata/out1.txt", 0, 0)

		require.Equal(t, nil, err, "Error occurred while copying")

		tmpBuff, err := os.ReadFile("testdata/out1.txt")

		require.Equal(t, nil, err, "Error occurred while reading created")

		os.Remove("testdata/out1.txt")

		require.Equal(t, []byte("hello"), tmpBuff, "Written not expected")
	})

	t.Run("Offset 1, limit 2", func(t *testing.T) {
		err := Copy("testdata/input1.txt", "testdata/out2.txt", 1, 2)

		require.Equal(t, nil, err, "Error occurred while copying")

		tmpBuff, err := os.ReadFile("testdata/out2.txt")

		require.Equal(t, nil, err, "Error occurred while reading created")

		os.Remove("testdata/out2.txt")

		require.Equal(t, []byte("el"), tmpBuff, "Written not expected")
	})

	t.Run("Offset 1, limit 10", func(t *testing.T) {
		err := Copy("testdata/input1.txt", "testdata/out3.txt", 1, 10)

		require.Equal(t, nil, err, "Error occurred while copying")

		tmpBuff, err := os.ReadFile("testdata/out3.txt")

		require.Equal(t, nil, err, "Error occurred while reading created")

		os.Remove("testdata/out3.txt")

		require.Equal(t, []byte("ello"), tmpBuff, "Written not expected")
	})
}
