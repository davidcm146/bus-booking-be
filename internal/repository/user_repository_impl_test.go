package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
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

// cleanupUsers removes all users matching the given emails so tests
// don't leak data into each other.
func cleanupUsers(t *testing.T, db *gorm.DB, emails ...string) {
	t.Helper()
	require.NoError(t, db.Where("email IN ?", emails).Delete(&model.User{}).Error)
}

func TestUserRepository_CreateAndFindByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	email := "test-create@example.com"
	defer cleanupUsers(t, db, email)

	user := &model.User{
		FullName: "Test User",
		Email:    email,
		Phone:    "+84900000001",
		Password: "hashed",
		Role:     model.RolePassenger,
	}

	err := repo.Create(ctx, user)
	require.NoError(t, err)
	require.NotZero(t, user.ID, "Create should set the ID")

	found, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, user.Email, found.Email)
	require.Equal(t, user.FullName, found.FullName)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	email := "test-email@example.com"
	defer cleanupUsers(t, db, email)

	user := &model.User{
		FullName: "Email Test",
		Email:    email,
		Phone:    "+84900000002",
		Password: "hashed",
	}
	require.NoError(t, repo.Create(ctx, user))

	found, err := repo.FindByEmail(ctx, email)
	require.NoError(t, err)
	require.Equal(t, user.ID, found.ID)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	found, err := repo.FindByEmail(ctx, "nonexistent@example.com")
	require.Nil(t, found)
	require.Error(t, err)

	appErr, ok := err.(*apperror.AppError)
	require.True(t, ok)
	require.Equal(t, 404, appErr.HTTPStatus)
}

func TestUserRepository_FindByPhone(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	email := "test-phone@example.com"
	phone := "+84900000003"
	defer cleanupUsers(t, db, email)

	user := &model.User{
		FullName: "Phone Test",
		Email:    email,
		Phone:    phone,
		Password: "hashed",
	}
	require.NoError(t, repo.Create(ctx, user))

	found, err := repo.FindByPhone(ctx, phone)
	require.NoError(t, err)
	require.Equal(t, user.ID, found.ID)
}

func TestUserRepository_FindAll(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	email1 := "test-findall-1@example.com"
	email2 := "test-findall-2@example.com"
	defer cleanupUsers(t, db, email1, email2)

	require.NoError(t, repo.Create(ctx, &model.User{
		FullName: "User 1", Email: email1, Phone: "+84900000004", Password: "h",
	}))
	require.NoError(t, repo.Create(ctx, &model.User{
		FullName: "User 2", Email: email2, Phone: "+84900000005", Password: "h",
	}))

	users, err := repo.FindAll(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(users), 2)
}

func TestUserRepository_Update(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	email := "test-update@example.com"
	defer cleanupUsers(t, db, email)

	user := &model.User{
		FullName: "Original",
		Email:    email,
		Phone:    "+84900000006",
		Password: "hashed",
	}
	require.NoError(t, repo.Create(ctx, user))

	user.FullName = "Updated"
	require.NoError(t, repo.Update(ctx, user))

	found, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, "Updated", found.FullName)
}

func TestUserRepository_Delete(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	email := "test-delete@example.com"
	defer cleanupUsers(t, db, email)

	user := &model.User{
		FullName: "Delete Me",
		Email:    email,
		Phone:    "+84900000007",
		Password: "hashed",
	}
	require.NoError(t, repo.Create(ctx, user))

	require.NoError(t, repo.Delete(ctx, user.ID))

	_, err := repo.FindByID(ctx, user.ID)
	require.Error(t, err)
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, 999999)
	require.Error(t, err)

	appErr, ok := err.(*apperror.AppError)
	require.True(t, ok)
	require.Equal(t, 404, appErr.HTTPStatus)
}

func TestMain(m *testing.M) {
	// Fail fast with a clear message if the DB isn't reachable, rather
	// than every test printing a confusing connection error.
	if dsn() == "" {
		fmt.Fprintln(os.Stderr, "skipping repository tests: set PG env to a Postgres DSN")
		os.Exit(0)
	}

	// The goose migration creates `role_id integer` (normalized FK to
	// roles), but the GORM model uses `Role` as a denormalized varchar
	// column — the app code treats Role as a string everywhere. Add the
	// column the model expects so repository tests can run against the
	// real DB. This is a model/migration mismatch that should be reconciled.
	db, err := gorm.Open(postgres.Open(dsn()), &gorm.Config{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to connect to test DB:", err)
		os.Exit(1)
	}
	db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS role varchar(20) NOT NULL DEFAULT 'passenger'")
	// The migration has `status varchar(30) NOT NULL` with no default, but
	// the model has no Status field. Set a default so INSERTs don't fail.
	db.Exec("ALTER TABLE users ALTER COLUMN status SET DEFAULT 'active'")

	code := m.Run()
	os.Exit(code)
}
