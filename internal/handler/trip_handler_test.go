package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/dto/response"
	"github.com/davidcm146/bus-booking-be/internal/mocks"
	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTripRouter(h *TripHandler) *gin.Engine {
	r := gin.New()
	v1 := r.Group("/api/v1")
	trips := v1.Group("/trips")
	trips.GET("", h.Search)
	trips.GET("/:id", h.GetByID)
	return r
}

func TestTripHandler_Search(t *testing.T) {
	valid := []struct {
		name  string
		query string
		want  request.SearchTripsRequest
	}{
		{
			"all filters",
			"?originId=28&destinationId=26&departureDate=2026-10-08&arrivalDate=2026-10-09",
			request.SearchTripsRequest{OriginID: 28, DestinationID: 26, DepartureDate: "2026-10-08", ArrivalDate: "2026-10-09"},
		},
		{"no filters", "", request.SearchTripsRequest{}},
	}
	for _, tt := range valid {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.MockTripService)
			r := newTripRouter(NewTripHandler(svc))

			svc.On("SearchTrips", mock.Anything, tt.want).Return([]response.TripResponse{{ID: 10}}, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/trips"+tt.query, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
			svc.AssertExpectations(t)
		})
	}

	invalid := []struct {
		name  string
		query string
	}{
		{"same origin and destination", "?originId=28&destinationId=28"},
		{"bad departure date format", "?departureDate=08/10/2026"},
		{"bad arrival date format", "?arrivalDate=08/10/2026"},
		{"negative area id", "?originId=-1&destinationId=26"},
		{"negative limit", "?limit=-5"},
		{"negative page", "?page=-1"},
	}
	for _, tt := range invalid {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.MockTripService)
			r := newTripRouter(NewTripHandler(svc))

			req := httptest.NewRequest(http.MethodGet, "/api/v1/trips"+tt.query, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
			svc.AssertNotCalled(t, "SearchTrips")
		})
	}
}

func TestTripHandler_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		setup      func(svc *mocks.MockTripService)
		wantStatus int
	}{
		{
			name: "success",
			path: "/api/v1/trips/10",
			setup: func(svc *mocks.MockTripService) {
				svc.On("GetTripByID", mock.Anything, 10).Return(&response.TripResponse{ID: 10}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "not found",
			path: "/api/v1/trips/999",
			setup: func(svc *mocks.MockTripService) {
				svc.On("GetTripByID", mock.Anything, 999).Return(nil, apperror.NotFound(i18n.KeyTripNotFound))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid id",
			path:       "/api/v1/trips/abc",
			setup:      func(svc *mocks.MockTripService) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.MockTripService)
			tt.setup(svc)
			r := newTripRouter(NewTripHandler(svc))

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))

			require.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
