package model

const (
	RolePassenger UserRole = "passenger"
	RoleOperator  UserRole = "operator"
	RoleAdmin     UserRole = "admin"
)

// Role is a row of the roles reference table.
type Role struct {
	ID   int      `gorm:"primaryKey" json:"id"`
	Name UserRole `json:"name"`
}
