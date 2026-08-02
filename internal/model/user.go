package model

import (
	"time"
)

// UserRole enumerates the allowed roles for a user.
type UserRole string

const (
	RolePassenger UserRole = "passenger"
	RoleOperator  UserRole = "operator"
	RoleAdmin     UserRole = "admin"
)

// User is the core domain entity representing a system user.
type User struct {
	ID        int        `gorm:"type:int;default:autoIncrement;primaryKey" json:"id"`
	Name      string     `gorm:"type:varchar(255);not null" json:"name"`
	Email     string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone     string     `gorm:"type:varchar(20)" json:"phone"`
	Password  string     `gorm:"type:varchar(255);not null" json:"-"`
	Role      UserRole   `gorm:"type:varchar(20);default:'passenger';not null" json:"role"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
