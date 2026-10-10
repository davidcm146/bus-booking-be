package request

// SearchTripsRequest is the inbound DTO for searching trips (query string).
// OriginID / DestinationID are administrative area ids (province or ward);
// DepartureDate is the local departure date in Asia/Ho_Chi_Minh.
type SearchTripsRequest struct {
	OriginID      int    `form:"originId" binding:"omitempty,gt=0"`
	DestinationID int    `form:"destinationId" binding:"omitempty,gt=0,nefield=OriginID"`
	DepartureDate string `form:"departureDate" binding:"omitempty,datetime=2006-01-02"`
	ArrivalDate   string `form:"arrivalDate" binding:"omitempty,datetime=2006-01-02"`
	Limit         int    `form:"limit" binding:"omitempty,gt=0"`
	Page          int    `form:"page" binding:"omitempty,gt=0"`
}
