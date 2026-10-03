package repository

import (
	"context"
	"testing"

	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/stretchr/testify/require"
)

func TestRoleRepository_FindByName(t *testing.T) {
	repo := NewRoleRepository(newTestDB(t))
	ctx := context.Background()

	t.Run("reference roles exist", func(t *testing.T) {
		for _, name := range []model.UserRole{model.RolePassenger, model.RoleOperator, model.RoleAdmin} {
			role, err := repo.FindByName(ctx, name)
			require.NoError(t, err, "roles reference data missing — run make migrate-up")
			require.Equal(t, name, role.Name)
		}
	})

	t.Run("unknown role", func(t *testing.T) {
		_, err := repo.FindByName(ctx, "pilot")
		requireAppError(t, err, 404, i18n.KeyRoleNotFound)
	})
}
