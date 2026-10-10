package service

import (
	"context"
	"time"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/dto/response"
	"github.com/davidcm146/bus-booking-be/internal/repository"
	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"github.com/davidcm146/bus-booking-be/internal/utils"
)

type tripServiceImpl struct {
	tripRepo repository.TripRepository
	now      func() time.Time
}

func NewTripService(tripRepo repository.TripRepository) TripService {
	return &tripServiceImpl{tripRepo: tripRepo, now: time.Now}
}

func (s *tripServiceImpl) SearchTrips(ctx context.Context, req request.SearchTripsRequest) ([]response.TripResponse, error) {
	now := s.now()
	limit, offset := utils.Paginate(req.Page, req.Limit, constant.DefaultTripPageSize, constant.MaxTripPageSize)
	filter := repository.TripSearchFilter{
		OriginAreaID:      req.OriginID,
		DestinationAreaID: req.DestinationID,
		DepartureFrom:     now,
		Limit:             limit,
		Offset:            offset,
	}

	if req.DepartureDate != "" {
		day, err := utils.ParseLocalDay(req.DepartureDate)
		if err != nil {
			return nil, err
		}
		filter.DepartureFrom, filter.DepartureTo = utils.Later(day, now), day.AddDate(0, 0, 1)
		if !filter.DepartureFrom.Before(filter.DepartureTo) {
			return []response.TripResponse{}, nil
		}
	}

	if req.ArrivalDate != "" {
		day, err := utils.ParseLocalDay(req.ArrivalDate)
		if err != nil {
			return nil, err
		}
		filter.ArrivalFrom, filter.ArrivalTo = day, day.AddDate(0, 0, 1)
		if !now.Before(filter.ArrivalTo) {
			return []response.TripResponse{}, nil
		}
	}

	trips, err := s.tripRepo.Search(ctx, filter)
	if err != nil {
		return nil, err
	}

	return response.ToTripResponses(trips), nil
}

func (s *tripServiceImpl) GetTripByID(ctx context.Context, id int) (*response.TripResponse, error) {
	trip, err := s.tripRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := response.ToTripResponse(trip)
	return &resp, nil
}
