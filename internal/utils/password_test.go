package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	t.Run("produces a bcrypt hash", func(t *testing.T) {
		hash, err := HashPassword("Str0ng!Pass")
		require.NoError(t, err)
		require.NotEmpty(t, hash)
		require.True(t, strings.HasPrefix(hash, "$2"))
	})

	t.Run("uses a random salt per call", func(t *testing.T) {
		h1, err := HashPassword("Str0ng!Pass")
		require.NoError(t, err)
		h2, err := HashPassword("Str0ng!Pass")
		require.NoError(t, err)
		require.NotEqual(t, h1, h2)
	})
}

func TestVerifyPassword(t *testing.T) {
	correctHash, err := HashPassword("Str0ng!Pass")
	require.NoError(t, err)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{"matching password", "Str0ng!Pass", correctHash, true},
		{"wrong password", "wrong", correctHash, false},
		{"empty password vs real hash", "", correctHash, false},
		{"empty vs empty hash", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, VerifyPassword(tt.password, tt.hash))
		})
	}
}
