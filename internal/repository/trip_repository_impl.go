package repository

import (
	"context"
	"errors"

	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"github.com/davidcm146/bus-booking-be/internal/shared/dberror"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"gorm.io/gorm"
)

// areaWithChildren selects an administrative area and its direct children.
const areaWithChildren = "SELECT id FROM administrative_areas WHERE (id = ? OR parent_id = ?) AND NOT deleted"

type tripRepositoryImpl struct {
	db *gorm.DB
}

// NewTripRepository creates a new GORM-backed TripRepository.
func NewTripRepository(db *gorm.DB) TripRepository {
	return &tripRepositoryImpl{db: db}
}

func (r *tripRepositoryImpl) active(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Where("trips.deleted = ?", false)
}

func (r *tripRepositoryImpl) Search(ctx context.Context, f TripSearchFilter) ([]model.Trip, error) {
	q := r.active(ctx).Where("trips.status = ?", model.TripStatusScheduled)
	if f.OriginAreaID > 0 {
		q = q.Where("trips.origin_id IN ("+areaWithChildren+")", f.OriginAreaID, f.OriginAreaID)
	}
	if f.DestinationAreaID > 0 {
		q = q.Where("trips.destination_id IN ("+areaWithChildren+")", f.DestinationAreaID, f.DestinationAreaID)
	}
	if !f.DepartureFrom.IsZero() {
		q = q.Where("trips.departure_time >= ?", f.DepartureFrom)
	}
	if !f.DepartureTo.IsZero() {
		q = q.Where("trips.departure_time < ?", f.DepartureTo)
	}
	if !f.ArrivalFrom.IsZero() {
		q = q.Where("trips.arrival_time >= ?", f.ArrivalFrom)
	}
	if !f.ArrivalTo.IsZero() {
		q = q.Where("trips.arrival_time < ?", f.ArrivalTo)
	}

	if f.Limit <= 0 {
		f.Limit = constant.DefaultTripPageSize
	}
	if f.Offset < 0 {
		f.Offset = 0
	}

	var trips []model.Trip
	err := q.Order("trips.departure_time, trips.id").Limit(f.Limit).Offset(f.Offset).Find(&trips).Error
	if err != nil {
		return nil, dberror.Translate(err, i18n.KeyTripFailedFetch)
	}
	return trips, nil
}

func (r *tripRepositoryImpl) FindByID(ctx context.Context, id int) (*model.Trip, error) {
	var trip model.Trip
	if err := r.active(ctx).First(&trip, "trips.id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound(i18n.KeyTripNotFound)
		}
		return nil, dberror.Translate(err, i18n.KeyTripFailedFind)
	}
	return &trip, nil
}
