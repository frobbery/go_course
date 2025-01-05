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

			"EMPTY": EnvValue{Value: ""},

			"FOO": EnvValue{Value: "   foo\nwith new line"},

			"HELLO": EnvValue{Value: "\"hello\""},

			"UNSET": EnvValue{NeedRemove: true},
		}

		require.Equal(t, environment, expectedMap, "Environment differs from expected")
	})
}
