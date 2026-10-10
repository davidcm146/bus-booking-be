package service

import (
	"context"
	"testing"
	"time"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/mocks"
	"github.com/davidcm146/bus-booking-be/internal/model"
	"github.com/davidcm146/bus-booking-be/internal/repository"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var businessLocation = time.FixedZone(constant.BusinessTimeZoneName, constant.BusinessTimeZoneOffset)

// fixedNow is 2026-10-06 09:15 in Asia/Ho_Chi_Minh.
var fixedNow = time.Date(2026, 10, 6, 9, 15, 0, 0, businessLocation)

func newTripSvc() (*tripServiceImpl, *mocks.MockTripRepository) {
	repo := new(mocks.MockTripRepository)
	return &tripServiceImpl{tripRepo: repo, now: func() time.Time { return fixedNow }}, repo
}

func TestSearchTrips(t *testing.T) {
	searchReq := func(departureDate, arrivalDate string) request.SearchTripsRequest {
		return request.SearchTripsRequest{OriginID: 28, DestinationID: 26, DepartureDate: departureDate, ArrivalDate: arrivalDate}
	}

	t.Run("future departure date searches the whole local day", func(t *testing.T) {
		svc, repo := newTripSvc()

		dayStart := time.Date(2026, 10, 8, 0, 0, 0, 0, businessLocation)
		repo.On("Search", mock.Anything, mock.MatchedBy(func(f repository.TripSearchFilter) bool {
			return f.OriginAreaID == 28 && f.DestinationAreaID == 26 &&
				f.DepartureFrom.Equal(dayStart) && f.DepartureTo.Equal(dayStart.AddDate(0, 0, 1)) &&
				f.ArrivalFrom.IsZero() && f.ArrivalTo.IsZero() &&
				f.Limit == constant.DefaultTripPageSize && f.Offset == 0
		})).Return([]model.Trip{{
			ID:                  10,
			BasePrice:           300000,
			Status:              model.TripStatusScheduled,
			AvailableSeatsCount: 40,
			SearchableSnapshot: model.TripSnapshot{
				Operator: model.SnapshotOperator{ID: 1, Code: "phuong-trang"},
				Origin:   model.SnapshotEndpoint{LocationID: 76, City: model.SnapshotCity{ID: 28}},
				Vehicle:  model.SnapshotVehicle{Type: "sleeper", Amenities: map[string]bool{"wifi": true}},
			},
		}}, nil)

		resp, err := svc.SearchTrips(context.Background(), searchReq("2026-10-08", ""))

		require.NoError(t, err)
		require.Len(t, resp, 1)
		require.Equal(t, 10, resp[0].ID)
		require.Equal(t, int64(300000), resp[0].BasePrice)
		require.Equal(t, 40, resp[0].AvailableSeats)
		require.Equal(t, "phuong-trang", resp[0].Operator.Code)
		require.Equal(t, 28, resp[0].Origin.City.ID)
		require.True(t, resp[0].Vehicle.Amenities["wifi"])
		repo.AssertExpectations(t)
	})

	t.Run("today only returns trips that have not departed", func(t *testing.T) {
		svc, repo := newTripSvc()

		repo.On("Search", mock.Anything, mock.MatchedBy(func(f repository.TripSearchFilter) bool {
			return f.DepartureFrom.Equal(fixedNow) &&
				f.DepartureTo.Equal(time.Date(2026, 10, 7, 0, 0, 0, 0, businessLocation))
		})).Return([]model.Trip{}, nil)

		resp, err := svc.SearchTrips(context.Background(), searchReq("2026-10-06", ""))

		require.NoError(t, err)
		require.Empty(t, resp)
		repo.AssertExpectations(t)
	})

	t.Run("arrival date searches the whole local day of upcoming trips", func(t *testing.T) {
		svc, repo := newTripSvc()

		dayStart := time.Date(2026, 10, 9, 0, 0, 0, 0, businessLocation)
		repo.On("Search", mock.Anything, mock.MatchedBy(func(f repository.TripSearchFilter) bool {
			return f.DepartureFrom.Equal(fixedNow) && f.DepartureTo.IsZero() &&
				f.ArrivalFrom.Equal(dayStart) && f.ArrivalTo.Equal(dayStart.AddDate(0, 0, 1))
		})).Return([]model.Trip{}, nil)

		resp, err := svc.SearchTrips(context.Background(), searchReq("", "2026-10-09"))

		require.NoError(t, err)
		require.Empty(t, resp)
		repo.AssertExpectations(t)
	})

	t.Run("departure and arrival dates are combined", func(t *testing.T) {
		svc, repo := newTripSvc()

		departureDay := time.Date(2026, 10, 8, 0, 0, 0, 0, businessLocation)
		arrivalDay := departureDay.AddDate(0, 0, 1)
		repo.On("Search", mock.Anything, mock.MatchedBy(func(f repository.TripSearchFilter) bool {
			return f.DepartureFrom.Equal(departureDay) && f.DepartureTo.Equal(arrivalDay) &&
				f.ArrivalFrom.Equal(arrivalDay) && f.ArrivalTo.Equal(arrivalDay.AddDate(0, 0, 1))
		})).Return([]model.Trip{}, nil)

		_, err := svc.SearchTrips(context.Background(), searchReq("2026-10-08", "2026-10-09"))

		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("no dates searches all upcoming trips", func(t *testing.T) {
		svc, repo := newTripSvc()

		repo.On("Search", mock.Anything, mock.MatchedBy(func(f repository.TripSearchFilter) bool {
			return f.DepartureFrom.Equal(fixedNow) && f.DepartureTo.IsZero() &&
				f.ArrivalFrom.IsZero() && f.ArrivalTo.IsZero()
		})).Return([]model.Trip{}, nil)

		_, err := svc.SearchTrips(context.Background(), request.SearchTripsRequest{})

		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	pagingCases := []struct {
		name        string
		limit, page int
		wantLimit   int
		wantOffset  int
	}{
		{"defaults", 0, 0, constant.DefaultTripPageSize, 0},
		{"page without limit uses default page size", 0, 2, constant.DefaultTripPageSize, constant.DefaultTripPageSize},
		{"limit and page", 10, 3, 10, 20},
		{"limit is capped", 500, 2, constant.MaxTripPageSize, constant.MaxTripPageSize},
	}
	for _, tt := range pagingCases {
		t.Run("paging: "+tt.name, func(t *testing.T) {
			svc, repo := newTripSvc()

			repo.On("Search", mock.Anything, mock.MatchedBy(func(f repository.TripSearchFilter) bool {
				return f.Limit == tt.wantLimit && f.Offset == tt.wantOffset
			})).Return([]model.Trip{}, nil)

			_, err := svc.SearchTrips(context.Background(), request.SearchTripsRequest{Limit: tt.limit, Page: tt.page})

			require.NoError(t, err)
			repo.AssertExpectations(t)
		})
	}

	pastCases := []struct {
		name string
		req  request.SearchTripsRequest
	}{
		{"past departure date", searchReq("2026-10-05", "")},
		{"past arrival date", searchReq("", "2026-10-05")},
	}
	for _, tt := range pastCases {
		t.Run(tt.name+" returns empty without querying", func(t *testing.T) {
			svc, repo := newTripSvc()

			resp, err := svc.SearchTrips(context.Background(), tt.req)

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Empty(t, resp)
			repo.AssertNotCalled(t, "Search")
		})
	}

	invalidCases := []struct {
		name string
		req  request.SearchTripsRequest
	}{
		{"invalid departure date", searchReq("06/10/2026", "")},
		{"invalid arrival date", searchReq("", "06/10/2026")},
	}
	for _, tt := range invalidCases {
		t.Run(tt.name+" is a bad request", func(t *testing.T) {
			svc, repo := newTripSvc()

			resp, err := svc.SearchTrips(context.Background(), tt.req)

			require.Nil(t, resp)
			appErr, ok := err.(*apperror.AppError)
			require.True(t, ok)
			require.Equal(t, 400, appErr.HTTPStatus)
			repo.AssertNotCalled(t, "Search")
		})
	}

	t.Run("repository error is returned", func(t *testing.T) {
		svc, repo := newTripSvc()

		repo.On("Search", mock.Anything, mock.Anything).Return(nil, apperror.Internal(i18n.KeyTripFailedFetch))

		resp, err := svc.SearchTrips(context.Background(), searchReq("2026-10-08", ""))

		require.Nil(t, resp)
		appErr, ok := err.(*apperror.AppError)
		require.True(t, ok)
		require.Equal(t, 500, appErr.HTTPStatus)
	})
}

func TestGetTripByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, repo := newTripSvc()

		repo.On("FindByID", mock.Anything, 10).Return(&model.Trip{ID: 10, Status: model.TripStatusScheduled}, nil)

		resp, err := svc.GetTripByID(context.Background(), 10)

		require.NoError(t, err)
		require.Equal(t, 10, resp.ID)
		require.Equal(t, model.TripStatusScheduled, resp.Status)
	})

	t.Run("not found", func(t *testing.T) {
		svc, repo := newTripSvc()

		repo.On("FindByID", mock.Anything, 999).Return(nil, apperror.NotFound(i18n.KeyTripNotFound))

		resp, err := svc.GetTripByID(context.Background(), 999)

		require.Nil(t, resp)
		appErr, ok := err.(*apperror.AppError)
		require.True(t, ok)
		require.Equal(t, 404, appErr.HTTPStatus)
	})
}
