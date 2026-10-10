package repository

import (
	"context"
	"testing"
	"time"

	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// seededAreas returns two distinct provinces from the seeded administrative_areas,
// each with one of its wards.
func seededAreas(t *testing.T, db *gorm.DB) (provinceA, wardA, provinceB, wardB int) {
	t.Helper()
	var rows []struct{ ProvinceID, WardID int }
	require.NoError(t, db.Raw(`
		SELECT DISTINCT ON (p.id) p.id AS province_id, w.id AS ward_id
		FROM administrative_areas p
		JOIN administrative_areas w ON w.parent_id = p.id AND NOT w.deleted
		WHERE p.parent_id IS NULL AND NOT p.deleted
		ORDER BY p.id, w.id
		LIMIT 2`).Scan(&rows).Error)
	require.Len(t, rows, 2, "administrative_areas seed missing — run the areas seed first")
	return rows[0].ProvinceID, rows[0].WardID, rows[1].ProvinceID, rows[1].WardID
}

func newTrip(origin, destination int, departure time.Time) *model.Trip {
	arrival := departure.Add(7 * time.Hour)
	return &model.Trip{
		OperatorID:          1,
		RouteID:             1,
		VehicleID:           1,
		OriginID:            origin,
		DestinationID:       destination,
		DepartureTime:       departure,
		ArrivalTime:         &arrival,
		BasePrice:           300000,
		Status:              model.TripStatusScheduled,
		AvailableSeatsCount: 40,
		SearchableSnapshot: model.TripSnapshot{
			Operator: model.SnapshotOperator{ID: 1, Code: "test-operator", Name: model.LocalizedText{"vi": "Nhà xe", "en": "Operator"}},
			Vehicle:  model.SnapshotVehicle{Type: "sleeper", TotalSeats: 40, Amenities: map[string]bool{"wifi": true}},
		},
	}
}

func createTrips(t *testing.T, db *gorm.DB, trips ...*model.Trip) {
	t.Helper()
	for _, trip := range trips {
		require.NoError(t, db.Create(trip).Error)
	}
	t.Cleanup(func() {
		ids := make([]int, len(trips))
		for i, trip := range trips {
			ids[i] = trip.ID
		}
		require.NoError(t, db.Delete(&model.Trip{}, ids).Error)
	})
}

func TestTripRepository_Search(t *testing.T) {
	db := newTestDB(t)
	repo := NewTripRepository(db)
	ctx := context.Background()
	provinceA, wardA, provinceB, wardB := seededAreas(t, db)

	// Far-future day so seeded trips never fall into the window.
	day := time.Date(2099, 1, 15, 0, 0, 0, 0, time.FixedZone("ICT", 7*60*60))
	late := newTrip(wardA, wardB, day.Add(20*time.Hour))
	early := newTrip(wardA, wardB, day.Add(6*time.Hour))
	reverse := newTrip(wardB, wardA, day.Add(8*time.Hour))
	nextDay := newTrip(wardA, wardB, day.Add(30*time.Hour))
	cancelled := newTrip(wardA, wardB, day.Add(9*time.Hour))
	cancelled.Status = "cancelled"
	deleted := newTrip(wardA, wardB, day.Add(10*time.Hour))
	deleted.Deleted = true
	createTrips(t, db, late, early, reverse, nextDay, cancelled, deleted)

	nextDayStart := day.AddDate(0, 0, 1)
	tests := []struct {
		name    string
		filter  TripSearchFilter
		wantIDs []int
	}{
		{"province matches trips from its wards, ordered by departure", TripSearchFilter{OriginAreaID: provinceA, DestinationAreaID: provinceB, DepartureFrom: day, DepartureTo: nextDayStart, Limit: 10}, []int{early.ID, late.ID}},
		{"ward ids match directly", TripSearchFilter{OriginAreaID: wardA, DestinationAreaID: wardB, DepartureFrom: day, DepartureTo: nextDayStart, Limit: 10}, []int{early.ID, late.ID}},
		{"direction matters", TripSearchFilter{OriginAreaID: provinceB, DestinationAreaID: provinceA, DepartureFrom: day, DepartureTo: nextDayStart, Limit: 10}, []int{reverse.ID}},
		{"limit is applied", TripSearchFilter{OriginAreaID: provinceA, DestinationAreaID: provinceB, DepartureFrom: day, DepartureTo: nextDayStart, Limit: 1}, []int{early.ID}},
		{"arrival window filters on arrival time", TripSearchFilter{OriginAreaID: provinceA, DestinationAreaID: provinceB, DepartureFrom: day, ArrivalFrom: nextDayStart, ArrivalTo: nextDayStart.AddDate(0, 0, 1), Limit: 10}, []int{late.ID, nextDay.ID}},
		{"zero areas match any area", TripSearchFilter{DepartureFrom: day, DepartureTo: nextDayStart.AddDate(0, 0, 1), Limit: 10}, []int{early.ID, reverse.ID, late.ID, nextDay.ID}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trips, err := repo.Search(ctx, tt.filter)
			require.NoError(t, err)

			ids := make([]int, len(trips))
			for i, trip := range trips {
				ids[i] = trip.ID
			}
			require.Equal(t, tt.wantIDs, ids)
		})
	}

	t.Run("snapshot round-trips through jsonb", func(t *testing.T) {
		trips, err := repo.Search(ctx, TripSearchFilter{
			OriginAreaID: wardA, DestinationAreaID: wardB,
			DepartureFrom: day, DepartureTo: day.AddDate(0, 0, 1), Limit: 1,
		})
		require.NoError(t, err)
		require.Len(t, trips, 1)
		require.Equal(t, "test-operator", trips[0].SearchableSnapshot.Operator.Code)
		require.Equal(t, "Nhà xe", trips[0].SearchableSnapshot.Operator.Name["vi"])
		require.True(t, trips[0].SearchableSnapshot.Vehicle.Amenities["wifi"])
	})
}

func TestTripRepository_FindByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewTripRepository(db)
	ctx := context.Background()
	_, wardA, _, wardB := seededAreas(t, db)

	departure := time.Date(2099, 2, 1, 8, 0, 0, 0, time.UTC)
	trip := newTrip(wardA, wardB, departure)
	deleted := newTrip(wardA, wardB, departure)
	deleted.Deleted = true
	createTrips(t, db, trip, deleted)

	t.Run("found", func(t *testing.T) {
		found, err := repo.FindByID(ctx, trip.ID)
		require.NoError(t, err)
		require.Equal(t, trip.ID, found.ID)
		require.True(t, found.DepartureTime.Equal(departure))
	})

	t.Run("deleted trip is not found", func(t *testing.T) {
		_, err := repo.FindByID(ctx, deleted.ID)
		requireAppError(t, err, 404, i18n.KeyTripNotFound)
	})

	t.Run("missing trip is not found", func(t *testing.T) {
		_, err := repo.FindByID(ctx, 999999999)
		requireAppError(t, err, 404, i18n.KeyTripNotFound)
	})
}
