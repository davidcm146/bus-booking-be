package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// dsn returns the test database connection string. It reads the PG
// environment variable (matching bun's convention) and falls back to
// the local docker-compose Postgres.
func dsn() string {
	if v := os.Getenv("PG"); v != "" {
		return v
	}
	return "postgres://postgres:10042003@localhost:5432/bus_booking?sslmode=disable"
}

// newTestDB opens a GORM connection to the test database. The caller
// is responsible for cleaning up rows it creates.
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(postgres.Open(dsn()), &gorm.Config{})
	require.NoError(t, err)
	return db
}

// cleanupUsers hard-deletes all users with the given phones once the test
// finishes so tests don't leak data into each other.
func cleanupUsers(t *testing.T, db *gorm.DB, phones ...string) {
	t.Helper()
	t.Cleanup(func() {
		require.NoError(t, db.Where("phone IN ?", phones).Delete(&model.User{}).Error)
	})
}

func newUser(name, phone string, email *string) *model.User {
	return &model.User{
		FullName: name,
		Email:    email,
		Phone:    &phone,
		Password: "hashed",
		Status:   model.UserStatusActive,
	}
}

func requireAppError(t *testing.T, err error, status int, key constant.MsgKey) {
	t.Helper()
	appErr, ok := err.(*apperror.AppError)
	require.True(t, ok, "expected *apperror.AppError, got %T", err)
	require.Equal(t, status, appErr.HTTPStatus)
	require.Equal(t, key, appErr.MessageKey)
}

func TestUserRepository_CreateAndFindByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	phone := "+84911000001"
	cleanupUsers(t, db, phone)

	role, err := NewRoleRepository(db).FindByName(ctx, model.RolePassenger)
	require.NoError(t, err, "roles reference data missing — run make migrate-up")

	user := newUser("Test User", phone, new("test-create@example.com"))
	user.RoleID = &role.ID
	require.NoError(t, repo.Create(ctx, user))
	require.NotZero(t, user.ID, "Create should set the ID")
	require.Equal(t, "vi", user.Language, "Create should return the DB default language")

	found, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, "Test User", found.FullName)
	require.Equal(t, "test-create@example.com", *found.Email)
	require.Equal(t, model.UserStatusActive, found.Status)
	require.Equal(t, model.RolePassenger, found.RoleName(), "FindByID should load the role")
}

func TestUserRepository_Create_WithoutEmail(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	phone1, phone2 := "+84911000002", "+84911000003"
	cleanupUsers(t, db, phone1, phone2)

	// Phone-only signups must store NULL, not '', so they don't collide on the email unique key.
	require.NoError(t, repo.Create(ctx, newUser("", phone1, nil)))
	require.NoError(t, repo.Create(ctx, newUser("", phone2, nil)))

	found, err := repo.FindByPhone(ctx, phone2)
	require.NoError(t, err)
	require.Nil(t, found.Email)
}

func TestUserRepository_Create_Conflict(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	phone, otherPhone := "+84911000004", "+84911000005"
	email := "test-conflict@example.com"
	cleanupUsers(t, db, phone, otherPhone)

	require.NoError(t, repo.Create(ctx, newUser("Existing", phone, &email)))

	tests := []struct {
		name string
		user *model.User
		key  constant.MsgKey
	}{
		{"duplicate phone", newUser("Dup", phone, new("other-conflict@example.com")), i18n.KeyUserPhoneConflict},
		{"duplicate email", newUser("Dup", otherPhone, &email), i18n.KeyUserEmailConflict},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireAppError(t, repo.Create(ctx, tt.user), 409, tt.key)
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	phone := "+84911000006"
	email := "test-email@example.com"
	cleanupUsers(t, db, phone)

	user := newUser("Email Test", phone, &email)
	require.NoError(t, repo.Create(ctx, user))

	found, err := repo.FindByEmail(ctx, email)
	require.NoError(t, err)
	require.Equal(t, user.ID, found.ID)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	repo := NewUserRepository(newTestDB(t))

	found, err := repo.FindByEmail(context.Background(), "nonexistent@example.com")
	require.Nil(t, found)
	requireAppError(t, err, 404, i18n.KeyUserNotFound)
}

func TestUserRepository_FindByPhone(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	phone := "+84911000007"
	cleanupUsers(t, db, phone)

	user := newUser("Phone Test", phone, nil)
	require.NoError(t, repo.Create(ctx, user))

	found, err := repo.FindByPhone(ctx, phone)
	require.NoError(t, err)
	require.Equal(t, user.ID, found.ID)
}

func TestUserRepository_FindAll_ExcludesDeleted(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	phone1, phone2 := "+84911000008", "+84911000009"
	cleanupUsers(t, db, phone1, phone2)

	kept, gone := newUser("Kept", phone1, nil), newUser("Gone", phone2, nil)
	require.NoError(t, repo.Create(ctx, kept))
	require.NoError(t, repo.Create(ctx, gone))
	require.NoError(t, repo.Delete(ctx, gone.ID))

	users, err := repo.FindAll(ctx)
	require.NoError(t, err)

	ids := make([]int, len(users))
	for i, u := range users {
		ids[i] = u.ID
	}
	require.Contains(t, ids, kept.ID)
	require.NotContains(t, ids, gone.ID)
}

func TestUserRepository_Update(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	phone, otherPhone := "+84911000010", "+84911000011"
	cleanupUsers(t, db, phone, otherPhone)

	user := newUser("Original", phone, new("test-update@example.com"))
	other := newUser("Other", otherPhone, new("test-update-other@example.com"))
	require.NoError(t, repo.Create(ctx, user))
	require.NoError(t, repo.Create(ctx, other))

	t.Run("updates fields", func(t *testing.T) {
		user.FullName = "Updated"
		require.NoError(t, repo.Update(ctx, user))

		found, err := repo.FindByID(ctx, user.ID)
		require.NoError(t, err)
		require.Equal(t, "Updated", found.FullName)
	})

	t.Run("email taken by another user is a conflict", func(t *testing.T) {
		dup := *other
		dup.Email = user.Email
		requireAppError(t, repo.Update(ctx, &dup), 409, i18n.KeyUserEmailConflict)
	})

	t.Run("deleted user is not found", func(t *testing.T) {
		require.NoError(t, repo.Delete(ctx, other.ID))
		other.FullName = "Ghost"
		requireAppError(t, repo.Update(ctx, other), 404, i18n.KeyUserNotFound)
	})
}

func TestUserRepository_Delete_IsSoft(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	phone := "+84911000012"
	cleanupUsers(t, db, phone)

	user := newUser("Delete Me", phone, nil)
	require.NoError(t, repo.Create(ctx, user))
	require.NoError(t, repo.Delete(ctx, user.ID))

	_, err := repo.FindByID(ctx, user.ID)
	requireAppError(t, err, 404, i18n.KeyUserNotFound)

	var row model.User
	require.NoError(t, db.First(&row, "id = ?", user.ID).Error, "row should still exist")
	require.True(t, row.Deleted)
	require.NotNil(t, row.DeletedAt)

	requireAppError(t, repo.Delete(ctx, user.ID), 404, i18n.KeyUserNotFound)
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	repo := NewUserRepository(newTestDB(t))

	requireAppError(t, repo.Delete(context.Background(), 999999), 404, i18n.KeyUserNotFound)
}

func TestMain(m *testing.M) {
	// Fail fast with a clear message if the DB isn't reachable, rather
	// than every test printing a confusing connection error.
	if dsn() == "" {
		fmt.Fprintln(os.Stderr, "skipping repository tests: set PG env to a Postgres DSN")
		os.Exit(0)
	}

	os.Exit(m.Run())
}
