package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.

type (
	User struct {
		ID string `json:"id" validate:"len:36"`

		Name string

		Age int `validate:"min:18|max:50"`

		Email string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`

		Role UserRole `validate:"in:admin,stuff"`

		Phones []string `validate:"len:11"`

		meta json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header []byte

		Payload []byte

		Signature []byte
	}

	Response struct {
		Code int `validate:"in:200,404,500"`

		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in interface{}

		expectedErr error
	}{
		{
			App{"123456"}, errors.New("errors:\nVersion:123456 len not 5"),
		},

		{
			Token{[]byte{0}, []byte{1}, []byte{2}}, nil,
		},

		{
			App{"12345"}, nil,
		},

		{
			Response{201, ""}, errors.New("errors:\nCode:201 not in 200,404,500"),
		},

		{
			Response{200, ""}, nil,
		},

		{
			User{
				Age: 12,

				Role: "UserRole",

				Phones: []string{"12345678901", "1"},
			},
			//nolint:lll
			errors.New("errors:\nID: len not 36\nAge:12 lesser than 18\nEmail: not matches ^\\w+@\\w+\\.\\w+$\nRole:UserRole not in admin,stuff\nPhones[1]:1 len not 11"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt

			t.Parallel()

			err := Validate(tt.in)

			require.Equal(t, tt.expectedErr, err, "Error not expected")
		})
	}
}
