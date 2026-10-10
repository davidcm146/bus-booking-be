package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// TripStatus enumerates the lifecycle states of a trip.
type TripStatus string

const (
	TripStatusScheduled TripStatus = "scheduled"
)

type Trip struct {
	ID                  int `gorm:"primaryKey"`
	OperatorID          int
	RouteID             int
	VehicleID           int
	OriginID            int
	DestinationID       int
	DepartureTime       time.Time
	ArrivalTime         *time.Time
	BasePrice           int64
	Status              TripStatus
	AvailableSeatsCount int
	BookedSeatsCount    int
	HeldSeatsCount      int
	SearchableSnapshot  TripSnapshot
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Deleted             bool
	DeletedAt           *time.Time
}

// LocalizedText is a translated value keyed by locale, e.g. {"vi": "...", "en": "..."}.
type LocalizedText map[string]string

// TripSnapshot is the trips.searchable_snapshot JSONB document.
// The json tags describe the stored document's keys, not the API output.
type TripSnapshot struct {
	Operator    SnapshotOperator `json:"operator"`
	Origin      SnapshotEndpoint `json:"origin"`
	Destination SnapshotEndpoint `json:"destination"`
	Route       SnapshotRoute    `json:"route"`
	Vehicle     SnapshotVehicle  `json:"vehicle"`
}

type SnapshotOperator struct {
	ID      int           `json:"id"`
	Code    string        `json:"code"`
	Name    LocalizedText `json:"name"`
	Brand   *string       `json:"brand"`
	Hotline *string       `json:"hotline"`
	LogoURL *string       `json:"logo_url"`
}

type SnapshotCity struct {
	ID   int           `json:"id"`
	Name LocalizedText `json:"name"`
}

type SnapshotEndpoint struct {
	LocationID int           `json:"location_id"`
	Name       LocalizedText `json:"name"`
	City       SnapshotCity  `json:"city"`
	Address    *string       `json:"address"`
	Lat        *float64      `json:"lat"`
	Lng        *float64      `json:"lng"`
}

type SnapshotRoute struct {
	ID              int           `json:"id"`
	Name            LocalizedText `json:"name"`
	DistanceKm      float64       `json:"distance_km"`
	DurationMinutes int           `json:"duration_minutes"`
	StopsCount      int           `json:"stops_count"`
}

type SnapshotVehicle struct {
	ID         int             `json:"id"`
	Type       string          `json:"type"`
	TotalSeats int             `json:"total_seats"`
	Amenities  map[string]bool `json:"amenities"`
}

// Scan implements sql.Scanner for reading the JSONB column.
func (s *TripSnapshot) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	}
	return fmt.Errorf("unsupported type %T for TripSnapshot", src)
}

// Value implements driver.Valuer for writing the JSONB column.
func (s TripSnapshot) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return string(b), err
}
