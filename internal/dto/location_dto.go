package dto

// NearbyLocationsRequest 주변 위치 조회 요청
type NearbyLocationsRequest struct {
	Latitude  float64 `form:"latitude" binding:"required"`
	Longitude float64 `form:"longitude" binding:"required"`
	Radius    int     `form:"radius"`
	Type      string  `form:"type"`
}

// LocationResponse 위치 응답
type LocationResponse struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	Address        string  `json:"address"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Distance       float64 `json:"distance"`
	OperatingHours string  `json:"operating_hours,omitempty"`
	Contact        string  `json:"contact,omitempty"`
}
