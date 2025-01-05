package main

import (
	"os"
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Test command executor", func(t *testing.T) {
		os.Setenv("USER", "user")

		os.Setenv("UNSET", "unset")

		cmd := []string{"echo", "${USER}", "${UNSET}", "${BAR}"}

		environment := Environment{
			"BAR": EnvValue{Value: "bar"},

			"UNSET": EnvValue{NeedRemove: true},
		}

		_ = RunCmd(cmd, environment)

		var tmpBuff []byte

		_, err := os.Stdout.Read(tmpBuff)

		require.Equal(t, nil, err, "Error occurred while reading created")

		require.Equal(t, []byte("user bar"), tmpBuff, "Written not expected")
	})
}
