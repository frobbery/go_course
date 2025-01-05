package main

import (
	// "os".
	"testing"
)

func TestRunCmd(t *testing.T) {
	_ = t
	/*
		t.Run("Test command executor", func(t *testing.T) {
			os.Setenv("USER", "user")

			os.Setenv("UNSET", "unset")

			cmd := []string{"echo", "$USER", "$UNSET", "$BAR", ">", "testdata/test.txt"}

			environment := Environment{
				"BAR": EnvValue{Value: "bar"},

				"UNSET": EnvValue{NeedRemove: true},
			}

			_ = RunCmd(cmd, environment)

			tmpBuff, err := os.ReadFile("testdata/test.txt")

			require.Equal(t, nil, err, "Error occurred while reading created")

			os.Remove("testdata/test.txt")

			require.Equal(t, []byte("user bar"), tmpBuff, "Written not expected")
		})*/
}
