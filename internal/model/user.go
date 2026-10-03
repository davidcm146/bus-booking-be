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
	ID         int        `gorm:"primaryKey" json:"id"`
	Name       string     `json:"name"`
	Email      *string    `json:"email"`
	Phone      *string    `json:"phone"`
	RoleID     *int       `json:"role_id"`
	Role       *Role      `json:"role,omitempty"`
	OperatorID *int       `json:"operator_id"`
	Password   string     `json:"-"`
	Language   string     `gorm:"default:vi" json:"language"`
	Status     UserStatus `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Deleted    bool       `json:"-"`
	DeletedAt  *time.Time `json:"-"`
}

// RoleName returns the name of the loaded role, or "" when it was not loaded.
func (u *User) RoleName() UserRole {
	if u.Role == nil {
		return ""
	}
	return u.Role.Name
}
