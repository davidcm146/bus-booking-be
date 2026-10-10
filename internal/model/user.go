package model

import (
	"time"
)

// UserRole enumerates the allowed roles for a user.
type UserRole string

// UserStatus enumerates the lifecycle states of a user account.
type UserStatus string

const (
	UserStatusActive UserStatus = "active"
)

// User is the core domain entity representing a system user.
// Fields mirror the users table in migrations/; the migration is the source of truth.
type User struct {
	ID         int `gorm:"primaryKey"`
	FullName   string
	Email      *string
	Phone      *string
	RoleID     *int
	Role       *Role
	OperatorID *int
	Password   string `json:"-"`
	Language   string `gorm:"default:vi"`
	Status     UserStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Deleted    bool
	DeletedAt  *time.Time
}

// RoleName returns the name of the loaded role, or "" when it was not loaded.
func (u *User) RoleName() UserRole {
	if u.Role == nil {
		return ""
	}
	return u.Role.Name
}
