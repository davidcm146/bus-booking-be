package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key"

func TestGenerateAndParseToken(t *testing.T) {
	t.Run("round trip", func(t *testing.T) {
		token, err := GenerateToken(42, "passenger", testSecret, time.Hour)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		claims, err := ParseToken(token, testSecret)
		require.NoError(t, err)
		require.Equal(t, 42, claims.UserID)
		require.Equal(t, "passenger", claims.Role)
	})

	t.Run("expired token is rejected", func(t *testing.T) {
		token, err := GenerateToken(1, "admin", testSecret, -time.Hour)
		require.NoError(t, err)

		_, err = ParseToken(token, testSecret)
		require.Error(t, err)
	})

	t.Run("wrong secret is rejected", func(t *testing.T) {
		token, err := GenerateToken(1, "admin", testSecret, time.Hour)
		require.NoError(t, err)

		_, err = ParseToken(token, "different-secret")
		require.Error(t, err)
	})

	t.Run("garbage string is rejected", func(t *testing.T) {
		_, err := ParseToken("not-a-jwt", testSecret)
		require.Error(t, err)
	})
}
