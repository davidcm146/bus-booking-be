package response

import (
	"time"

	"github.com/davidcm146/bus-booking-be/internal/model"
)

// TripResponse is the outbound DTO for a trip search result or trip detail.
type TripResponse struct {
	ID             int                  `json:"id"`
	DepartureTime  time.Time            `json:"departureTime"`
	ArrivalTime    *time.Time           `json:"arrivalTime"`
	BasePrice      int64                `json:"basePrice"`
	Status         model.TripStatus     `json:"status"`
	AvailableSeats int                  `json:"availableSeats"`
	Operator       TripOperatorResponse `json:"operator"`
	Origin         TripEndpointResponse `json:"origin"`
	Destination    TripEndpointResponse `json:"destination"`
	Route          TripRouteResponse    `json:"route"`
	Vehicle        TripVehicleResponse  `json:"vehicle"`
}

type TripOperatorResponse struct {
	ID      int                 `json:"id"`
	Code    string              `json:"code"`
	Name    model.LocalizedText `json:"name"`
	Brand   *string             `json:"brand"`
	Hotline *string             `json:"hotline"`
	LogoURL *string             `json:"logoUrl"`
}

type TripCityResponse struct {
	ID   int                 `json:"id"`
	Name model.LocalizedText `json:"name"`
}

type TripEndpointResponse struct {
	LocationID int                 `json:"locationId"`
	Name       model.LocalizedText `json:"name"`
	City       TripCityResponse    `json:"city"`
	Address    *string             `json:"address"`
	Lat        *float64            `json:"lat"`
	Lng        *float64            `json:"lng"`
}

type TripRouteResponse struct {
	ID              int                 `json:"id"`
	Name            model.LocalizedText `json:"name"`
	DistanceKm      float64             `json:"distanceKm"`
	DurationMinutes int                 `json:"durationMinutes"`
	StopsCount      int                 `json:"stopsCount"`
}

type TripVehicleResponse struct {
	ID         int             `json:"id"`
	Type       string          `json:"type"`
	TotalSeats int             `json:"totalSeats"`
	Amenities  map[string]bool `json:"amenities"`
}

// ToTripResponse maps a domain model to a response DTO.
func ToTripResponse(t *model.Trip) TripResponse {
	s := t.SearchableSnapshot
	return TripResponse{
		ID:             t.ID,
		DepartureTime:  t.DepartureTime,
		ArrivalTime:    t.ArrivalTime,
		BasePrice:      t.BasePrice,
		Status:         t.Status,
		AvailableSeats: t.AvailableSeatsCount,
		Operator: TripOperatorResponse{
			ID:      s.Operator.ID,
			Code:    s.Operator.Code,
			Name:    s.Operator.Name,
			Brand:   s.Operator.Brand,
			Hotline: s.Operator.Hotline,
			LogoURL: s.Operator.LogoURL,
		},
		Origin:      toTripEndpointResponse(s.Origin),
		Destination: toTripEndpointResponse(s.Destination),
		Route: TripRouteResponse{
			ID:              s.Route.ID,
			Name:            s.Route.Name,
			DistanceKm:      s.Route.DistanceKm,
			DurationMinutes: s.Route.DurationMinutes,
			StopsCount:      s.Route.StopsCount,
		},
		Vehicle: TripVehicleResponse{
			ID:         s.Vehicle.ID,
			Type:       s.Vehicle.Type,
			TotalSeats: s.Vehicle.TotalSeats,
			Amenities:  s.Vehicle.Amenities,
		},
	}
}

// ToTripResponses maps a slice of domain models to response DTOs.
func ToTripResponses(trips []model.Trip) []TripResponse {
	responses := make([]TripResponse, len(trips))
	for i := range trips {
		responses[i] = ToTripResponse(&trips[i])
	}
	return responses
}

func toTripEndpointResponse(e model.SnapshotEndpoint) TripEndpointResponse {
	return TripEndpointResponse{
		LocationID: e.LocationID,
		Name:       e.Name,
		City:       TripCityResponse{ID: e.City.ID, Name: e.City.Name},
		Address:    e.Address,
		Lat:        e.Lat,
		Lng:        e.Lng,
	}
}
