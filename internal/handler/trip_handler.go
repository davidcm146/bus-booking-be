package handler

import (
	"strconv"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/service"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/davidcm146/bus-booking-be/internal/shared/response"
	"github.com/davidcm146/bus-booking-be/internal/shared/validator"
	"github.com/gin-gonic/gin"
)

// TripHandler handles HTTP requests for trip queries.
type TripHandler struct {
	tripService service.TripService
}

// NewTripHandler creates a new TripHandler with the given service.
func NewTripHandler(tripService service.TripService) *TripHandler {
	return &TripHandler{tripService: tripService}
}

// Search godoc
// @Summary      Search trips
// @Description  Upcoming scheduled trips, optionally filtered by origin/destination administrative area (province or ward) and local departure/arrival date (Asia/Ho_Chi_Minh), ordered by departure time.
// @Tags         trips
// @Produce      json
// @Param        originId       query  int     false  "Origin administrative area ID"
// @Param        destinationId  query  int     false  "Destination administrative area ID"
// @Param        departureDate  query  string  false  "Departure date (YYYY-MM-DD)"
// @Param        arrivalDate    query  string  false  "Arrival date (YYYY-MM-DD)"
// @Success      200  {array}  response.TripResponse
// @Failure      422
// @Router       /api/v1/trips [get]
func (h *TripHandler) Search(c *gin.Context) {
	var req request.SearchTripsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, validator.FormatError(c.Request.Context(), err))
		return
	}

	result, err := h.tripService.SearchTrips(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// GetByID godoc
// @Summary      Get trip by ID
// @Tags         trips
// @Produce      json
// @Param        id  path  int  true  "Trip ID"
// @Success      200  {object}  response.TripResponse
// @Failure      404
// @Failure      422
// @Router       /api/v1/trips/{id} [get]
func (h *TripHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ValidationError(c, i18n.T(c.Request.Context(), string(i18n.KeyErrorInvalidIDFormat), nil))
		return
	}

	result, err := h.tripService.GetTripByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}
