package main

import (
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Test environment extraction", func(t *testing.T) {
		environment, err := ReadDir("testdata/env")

		require.Equal(t, err, nil, "Error not nil")

		expectedMap := Environment{
			"BAR": EnvValue{Value: "bar"},

			"EMPTY": EnvValue{NeedRemove: true},

			"FOO": EnvValue{Value: "foo\x00with new line"},

			"HELLO": EnvValue{Value: "\"hello\""},

			"UNSET": EnvValue{NeedRemove: true},
		}

		require.Equal(t, environment, expectedMap, "Environment differs from expected")
	})
}
